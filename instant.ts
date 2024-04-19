'use strict'

import * as commandLineArgs from 'command-line-args'
import * as glob from 'glob'
import * as fs from 'fs'
import * as child from 'child_process'
import * as util from 'util'
import * as path from 'path'
import { env } from 'process'
import Ajv from 'ajv'

const schema = require('./schema/package-metadata.schema.json')
const ajv = new Ajv()
const validate = ajv.compile(schema)

const exec = util.promisify(child.exec)

interface PackageInfo {
  metadata: {
    id: string
    name: string
    description: string
    version: string
    dependencies: string[]
    environmentVariables: object
  }
  path: string
}

interface PackagesMap {
  [packageID: string]: PackageInfo
}

type EnvironmentVar = {
  'Environment Variable': string
  'Current Value': string | undefined
}

function getInstantOHIEPackages(): PackagesMap {
  const packages: PackagesMap = {}
  let metaPathRegex = 'package-metadata.json'
  let pathRegex = 'instant.json' // Keeping the instant.json logic to ensure backward compatibility
  let paths = [] as string[]
  let nestingLevel = 0

  while (nestingLevel < 5) {
    metaPathRegex = '*/' + metaPathRegex
    pathRegex = '*/' + pathRegex
    paths = paths.concat(glob.sync(metaPathRegex), glob.sync(pathRegex))
    nestingLevel += 1
  }

  for (const path of paths) {
    let metadata
    try {
      metadata = JSON.parse(fs.readFileSync(path).toString())
    } catch (err) {
      console.error(`❌ Failed to parse package metadata for ${path}.`)
      throw err
    }

    const isValid = validate(metadata)

    if (!isValid) {
      console.error(
        `❌ Package metadata for ${metadata.id} is invalid: ${(
          validate.errors || []
        )
          .map((error) => error.message)
          .join(', ')}`
      )
      throw new Error(`Invalid package metadata for ${metadata.id}`)
    }

    packages[metadata.id] = {
      metadata,
      path:
        path.includes('instant.json') === true
          ? path.replace('instant.json', '')
          : path.replace('package-metadata.json', '')
    }
  }

  return packages
}

let error = false

async function runBashScript(path: string, filename: string, args: string[]) {
  const cmd = `bash ${path}${filename} ${args.join(' ')}`

  try {
    const promise = exec(cmd)
    if (promise.child) {
      promise.child.stdout?.on('data', (data) => console.log('\t' + data))
      promise.child.stderr?.on('data', (data) => console.error('\t' + data))
    }
    await promise
  } catch (err) {
    console.error(`❌ Script ${path}${filename} returned an error`)
    error = true
  }
}

async function runTests(path: string) {
  const cmd = `node_modules/.bin/cucumber-js ${path}`

  try {
    const promise = exec(cmd)
    if (promise.child) {
      promise.child.stdout?.on('data', (data) => console.log(data))
      promise.child.stderr?.on('data', (data) => console.error(data))
    }
    await promise
  } catch (err) {
    console.error(`Error: Tests at ${path} returned an error`)
    console.log(err.stdout)
    console.log(err.stderr)
  }
}

export const createDependencyTree = (allPackages, chosenPackageIds) => {
  const tree = {}
  const visited = new Set()

  const addDependencies = (id, node) => {
    if (visited.has(id)) {
      throw new Error(
        `Circular dependency detected: ${id} has already been visited.`
      )
    }
    if (!allPackages[id] || !allPackages[id].metadata) {
      throw new Error(`Invalid package ID: ${id}`)
    }

    visited.add(id)
    const deps = allPackages[id].metadata.dependencies || []
    deps.forEach((dep) => {
      if (!node[dep]) {
        node[dep] = {}
        addDependencies(dep, node[dep])
      }
    })
    visited.delete(id)
  }

  chosenPackageIds.forEach((id) => {
    if (!tree[id]) {
      tree[id] = {}
      addDependencies(id, tree[id])
    }
  })
  return tree
}

export const walkDependencyTree = async (tree, preOrPost, action) => {
  const visitNode = async (node) => {
    await Promise.all(
      Object.keys(node).map(async (key) => {
        if (preOrPost === 'pre') await action(key)
        await visitNode(node[key])
        if (preOrPost === 'post') await action(key)
      })
    )
  }
  await visitNode(tree)
}

export const concurrentifyAction = (
  action: (id: string) => Promise<void>,
  maxConcurrentActions: number
) => {
  const idToPromiseMap: Map<string, Promise<void>> = new Map()
  const activePromises: Promise<void>[] = []

  const concurrentAction = async (
    id: string,
    action: (id: string) => Promise<void>
  ) => {
    // Wait until there's space to start a new action
    while (activePromises.length >= maxConcurrentActions) {
      await Promise.race(activePromises)
    }

    if (!idToPromiseMap.has(id)) {
      const promise = action(id).then(() => {
        // Remove itself from activePromises once done
        activePromises.splice(activePromises.indexOf(promise), 1)
      })
      activePromises.push(promise)
      idToPromiseMap.set(id, promise)
      return promise
    } else {
      return idToPromiseMap.get(id)
    }
  }

  return async (id: string) => {
    return concurrentAction(id, action)
  }
}

const truncateString = (str: string, maxLength: number) => {
  return str.length > maxLength
    ? `${str.substring(0, maxLength - 10)}...[trunc]`
    : str
}

const setEnvVars = (packageInfo: PackageInfo) => {
  const envVars = [] as EnvironmentVar[]

  for (let envVar in packageInfo.metadata.environmentVariables) {
    const defaultEnv = packageInfo.metadata.environmentVariables[envVar]
    if (env[envVar] === undefined || env[envVar] === null) {
      process.env[envVar] = defaultEnv
    }

    envVars.push({
      'Environment Variable': envVar,
      'Current Value': env[envVar]
    })
  }

  if (envVars?.length > 0) {
    console.log(
      `🛠️ Config set for ${packageInfo.metadata.name} (${packageInfo.metadata.id}):`
    )
    console.table(
      envVars.map(
        ({ 'Environment Variable': envVar, 'Current Value': currVal }) => ({
          'Environment Variable': truncateString(envVar, 50),
          'Current Value': truncateString(currVal || '', 50)
        })
      )
    )
  }
}

// Main script execution
const main = async () => {
  const allPackages = getInstantOHIEPackages()
  console.log(
    `📦 Found ${Object.keys(allPackages).length} packages: ${Object.values(
      allPackages
    )
      .map((p) => p.metadata.id)
      .join(', ')}`
  )

  const main = commandLineArgs(
    [
      {
        name: 'command',
        defaultOption: true
      }
    ],
    {
      stopAtFirstUnknown: true
    }
  )

  let argv = main._unknown || []

  // main commands
  if (['init', 'up', 'down', 'destroy'].includes(main.command)) {
    const mainOptions = commandLineArgs(
      [
        {
          name: 'target',
          alias: 't',
          defaultValue: 'docker'
        },
        {
          name: 'only',
          alias: 'o',
          type: Boolean
        },
        {
          name: 'dev',
          alias: 'd',
          type: Boolean
        },
        {
          name: 'concurrency',
          alias: 'c',
          type: Number,
          defaultValue: 5
        }
      ],
      { argv, stopAtFirstUnknown: true }
    )

    console.log(`🎯 Target environment is: ${mainOptions.target}`)
    console.log(`🔀 Running using concurrency of ${mainOptions.concurrency}`)

    argv = mainOptions._unknown || []
    let chosenPackageIds = argv

    if (
      !chosenPackageIds.every((id) => Object.keys(allPackages).includes(id))
    ) {
      throw new Error(
        `Deploy - Unknown package id in list: ${chosenPackageIds}`
      )
    }

    if (chosenPackageIds.length < 1) {
      chosenPackageIds = Object.keys(allPackages)
    }

    if (mainOptions.dev) {
      mainOptions.mode = 'dev'
    } else {
      mainOptions.mode = 'prod'
    }

    console.log(
      `👉 Selected package IDs to operate on: ${chosenPackageIds.join(', ')}\n`
    )

    const dependencyTree = createDependencyTree(allPackages, chosenPackageIds)

    const action = async (id) => {
      switch (main.command) {
        case 'init':
          console.log(
            `🚀 Initializing package ${allPackages[id].metadata.name} (${id})...`
          )
          break
        case 'up':
          console.log(
            `🚀 Starting package ${allPackages[id].metadata.name} (${id})...`
          )
          break
        case 'down':
          console.log(
            `🛑 Stopping package ${allPackages[id].metadata.name} (${id})...`
          )
          break
        case 'destroy':
          console.log(
            `🧨 Destroying package ${allPackages[id].metadata.name} (${id})...`
          )
          break
        default:
          console.log(
            `🚀 Performing action on package ${allPackages[id].metadata.name} (${id})...`
          )
      }

      setEnvVars(allPackages[id])
      let scriptPath, scriptName
      switch (mainOptions.target) {
        case 'docker':
          scriptPath = `${allPackages[id].path}docker/`
          scriptName = 'compose.sh'
          break
        case 'k8s':
        case 'kubernetes':
          scriptPath = `${allPackages[id].path}kubernetes/main/`
          scriptName = 'k8s.sh'
          break
        default:
          scriptPath = `${allPackages[id].path}`
          scriptName = 'swarm.sh'
      }
      const scriptArgs =
        mainOptions.target === 'swarm'
          ? [main.command, mainOptions.mode]
          : [main.command]
      await runBashScript(scriptPath, scriptName, scriptArgs)
    }

    // execute action
    if (mainOptions.only) {
      for (const id of chosenPackageIds) {
        await action(id)
      }
    } else if (['destroy', 'down'].includes(main.command)) {
      await walkDependencyTree(
        dependencyTree,
        'pre',
        concurrentifyAction(action, mainOptions.concurrency)
      )
    } else {
      await walkDependencyTree(
        dependencyTree,
        'post',
        concurrentifyAction(action, mainOptions.concurrency)
      )
    }

    if (error) {
      console.log('\n❌ Some scripts returned errors')
    } else {
      console.log('\n🟢 Success!')
    }
  }

  // test command
  if (main.command === 'test') {
    const testOptions = commandLineArgs(
      [
        {
          name: 'host',
          alias: 'h',
          defaultValue: 'localhost'
        },
        {
          name: 'port',
          alias: 'p',
          defaultValue: '5000'
        }
      ],
      { argv, stopAtFirstUnknown: true }
    )

    argv = testOptions._unknown || []
    let chosenPackageIds = argv

    if (
      !chosenPackageIds.every((id) => Object.keys(allPackages).includes(id))
    ) {
      throw new Error(
        `Testing - Unknown package id in list: ${chosenPackageIds}`
      )
    }

    if (chosenPackageIds.length < 1) {
      chosenPackageIds = Object.keys(allPackages)
    }

    // Order the packages such that the dependencies are instantiated first
    const orderedIds: string[] = []
    const orderIdsAction = (id: string) => {
      if (!orderedIds.includes(id)) {
        orderedIds.push(id)
      }
    }

    const dependencyTree = createDependencyTree(allPackages, chosenPackageIds)
    walkDependencyTree(dependencyTree, 'post', orderIdsAction)
    chosenPackageIds = orderedIds

    console.log(`Running tests for packages: ${chosenPackageIds.join(', ')}`)
    console.log(`Using host: ${testOptions.host}:${testOptions.port}`)

    for (const id of chosenPackageIds) {
      const features = path.resolve(allPackages[id].path, 'features')
      await runTests(features)
    }
  }
}

if (process.env.NODE_ENV !== 'test') {
  // Entry point IIFE with base error handling
  ;(async () => {
    try {
      await main()
    } catch (error) {
      console.log(error)
      process.exit(1)
    }
  })()
}
