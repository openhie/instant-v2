'use strict'

import * as commandLineArgs from 'command-line-args'
import * as glob from 'glob'
import * as fs from 'fs'
import * as child from 'child_process'
import * as util from 'util'
import * as path from 'path'
import { env } from 'process'

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
  'Default Value': string
  'Updated Value': string | undefined
}

function getInstantOHIEPackages(): PackagesMap {
  const packages: PackagesMap = {}
  let metaPathRegex = 'package-metadata.json'
  let pathRegex = 'instant.json' //Keeping the instant.json logic to ensure backward compatibility
  let paths = [] as string[]
  let nestingLevel = 0

  while (nestingLevel < 5) {
    metaPathRegex = '*/' + metaPathRegex
    pathRegex = '*/' + pathRegex
    paths = paths.concat(glob.sync(metaPathRegex), glob.sync(pathRegex))
    nestingLevel += 1
  }

  for (const path of paths) {
    try {
      const metadata = JSON.parse(fs.readFileSync(path).toString())
      packages[metadata.id] = {
        metadata,
        path:
          path.includes('instant.json') === true
            ? path.replace('instant.json', '')
            : path.replace('package-metadata.json', '')
      }
    } catch (err) {
      console.error(`Failed to parse package metadata for ${path}.`)
      throw err
    }
  }

  return packages
}

async function runBashScript(path: string, filename: string, args: string[]) {
  const cmd = `bash ${path}${filename} ${args.join(' ')}`
  console.log(`Executing: ${cmd}`)

  try {
    const promise = exec(cmd)
    if (promise.child) {
      promise.child.stdout?.on('data', (data) => console.log(data))
      promise.child.stderr?.on('data', (data) => console.error(data))
    }
    await promise
  } catch (err) {
    console.error(`Script ${path}${filename} returned an error`)
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

const createDependencyTree = (allPackages, chosenPackageIds) => {
  const tree = {}
  const addDependencies = (id, node) => {
    if (!allPackages[id] || !allPackages[id].metadata) return
    const deps = allPackages[id].metadata.dependencies || []
    deps.forEach((dep) => {
      if (!node[dep]) {
        node[dep] = {}
        addDependencies(dep, node[dep])
      }
    })
  }
  chosenPackageIds.forEach((id) => {
    if (!tree[id]) {
      tree[id] = {}
      addDependencies(id, tree[id])
    }
  })
  return tree
}

const walkDependencyTree = async (tree, preOrPost, action) => {
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

const concurrentifyAction = (
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
      console.log(`Executing action for ${id}`)
      const promise = action(id).then(() => {
        console.log(`Done ${id}`)
        // Remove itself from activePromises once done
        activePromises.splice(activePromises.indexOf(promise), 1)
      })
      activePromises.push(promise)
      idToPromiseMap.set(id, promise)
      return promise
    } else {
      console.log(`Returning existing promise for ${id}`)
      return idToPromiseMap.get(id)
    }
  }

  return async (id: string) => {
    return concurrentAction(id, action)
  }
}

const setEnvVars = (packageInfo: PackageInfo) => {
  console.log(
    `------------------------------------------------------------\nConfig Details: ${packageInfo.metadata.name} (${packageInfo.metadata.id})\n------------------------------------------------------------`
  )
  const envVars = [] as EnvironmentVar[]

  for (let envVar in packageInfo.metadata.environmentVariables) {
    const defaultEnv = packageInfo.metadata.environmentVariables[envVar]
    if (env[envVar] === undefined || env[envVar] === null) {
      process.env[envVar] = defaultEnv
    }

    envVars.push({
      'Environment Variable': envVar,
      'Default Value': defaultEnv,
      'Updated Value': env[envVar]
    })
  }

  if (envVars?.length > 0) {
    console.table(envVars)
  }
}

// Main script execution
const main = async () => {
  const allPackages = getInstantOHIEPackages()
  console.log(
    `Found ${Object.keys(allPackages).length} packages: ${Object.values(
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
        }
      ],
      { argv, stopAtFirstUnknown: true }
    )

    console.log(`Target environment is: ${mainOptions.target}`)

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
      `Selected package IDs to operate on: ${chosenPackageIds.join(', ')}`
    )

    const dependencyTree = createDependencyTree(allPackages, chosenPackageIds)
    console.log(JSON.stringify(dependencyTree, null, 2))
    const CONCURRENT_ACTIONS = 10

    const action = async (id) => {
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
    console.log(mainOptions)
    if (mainOptions.only) {
      for (const id of chosenPackageIds) {
        await action(id)
      }
    } else if (['destroy', 'down'].includes(main.command)) {
      walkDependencyTree(
        dependencyTree,
        'pre',
        concurrentifyAction(action, CONCURRENT_ACTIONS)
      )
    } else {
      walkDependencyTree(
        dependencyTree,
        'post',
        concurrentifyAction(action, CONCURRENT_ACTIONS)
      )
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

// Entry point IIFE with base error handling
;(async () => {
  try {
    await main()
  } catch (error) {
    console.log(error)
    process.exit(1)
  }
})()
