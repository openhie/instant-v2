import {
  createDependencyTree,
  walkDependencyTree,
  concurrentifyAction
} from './instant'

describe('createDependencyTree', () => {
  it('should handle a single package with no dependencies', () => {
    const allPackages = {
      package1: { metadata: { id: 'package1', dependencies: [] } }
    }
    const chosenPackageIds = ['package1']
    const expectedTree = { package1: {} }
    expect(createDependencyTree(allPackages, chosenPackageIds)).toEqual(
      expectedTree
    )
  })

  it('should handle multiple packages with dependencies', () => {
    const allPackages = {
      package1: { metadata: { id: 'package1', dependencies: ['package2'] } },
      package2: { metadata: { id: 'package2', dependencies: [] } }
    }
    const chosenPackageIds = ['package1']
    const expectedTree = { package1: { package2: {} } }
    expect(createDependencyTree(allPackages, chosenPackageIds)).toEqual(
      expectedTree
    )
  })

  it('should handle complex dependency trees', () => {
    const allPackages = {
      package1: {
        metadata: { id: 'package1', dependencies: ['package2', 'package3'] }
      },
      package2: { metadata: { id: 'package2', dependencies: ['package3'] } },
      package3: { metadata: { id: 'package3', dependencies: [] } }
    }
    const chosenPackageIds = ['package1']
    const expectedTree = {
      package1: { package2: { package3: {} }, package3: {} }
    }
    expect(createDependencyTree(allPackages, chosenPackageIds)).toEqual(
      expectedTree
    )
  })

  it('should throw an error for circular dependencies', () => {
    const allPackages = {
      package1: { metadata: { id: 'package1', dependencies: ['package2'] } },
      package2: { metadata: { id: 'package2', dependencies: ['package1'] } }
    }
    const chosenPackageIds = ['package1']
    expect(() => createDependencyTree(allPackages, chosenPackageIds)).toThrow(
      'Circular dependency detected: package1 has already been visited.'
    )
  })

  it('should handle invalid or missing package IDs gracefully', () => {
    const allPackages = {
      package1: { metadata: { id: 'package1', dependencies: [] } }
    }
    const chosenPackageIds = ['nonExistentPackage']
    expect(() => createDependencyTree(allPackages, chosenPackageIds)).toThrow(
      'Invalid package ID: nonExistentPackage'
    )
  })
})

describe('walkDependencyTree', () => {
  const mockAction = jest.fn()

  beforeEach(() => {
    mockAction.mockClear()
  })

  it('should call action on a single node tree in pre-order', async () => {
    const tree = { package1: {} }
    await walkDependencyTree(tree, 'pre', mockAction)
    expect(mockAction).toHaveBeenCalledTimes(1)
    expect(mockAction).toHaveBeenCalledWith('package1')
  })

  it('should call action on a single node tree in post-order', async () => {
    const tree = { package1: {} }
    await walkDependencyTree(tree, 'post', mockAction)
    expect(mockAction).toHaveBeenCalledTimes(1)
    expect(mockAction).toHaveBeenCalledWith('package1')
  })

  it('should walk a complex tree in pre-order and call action correctly', async () => {
    const tree = {
      package1: {
        package2: {},
        package3: {
          package4: {}
        }
      }
    }
    await walkDependencyTree(tree, 'pre', mockAction)
    expect(mockAction.mock.calls).toEqual([
      ['package1'],
      ['package2'],
      ['package3'],
      ['package4']
    ])
  })

  it('should walk a complex tree in post-order and call action correctly', async () => {
    const tree = {
      package1: {
        package2: {},
        package3: {
          package4: {}
        }
      }
    }
    await walkDependencyTree(tree, 'post', mockAction)
    expect(mockAction.mock.calls).toEqual([
      ['package2'],
      ['package4'],
      ['package3'],
      ['package1']
    ])
  })

  it('should handle an empty tree', async () => {
    const tree = {}
    await walkDependencyTree(tree, 'pre', mockAction)
    expect(mockAction).not.toHaveBeenCalled()
  })
})

describe('concurrentifyAction', () => {
  it('executes actions concurrently up to the specified limit', async () => {
    const action = jest
      .fn()
      .mockImplementation(
        (id) => new Promise((resolve) => setTimeout(resolve, 100))
      )
    const concurrentAction = concurrentifyAction(action, 2)

    const startTime = Date.now()
    await Promise.all([
      concurrentAction('1'),
      concurrentAction('2'),
      concurrentAction('3') // This should wait until one of the first two completes
    ])
    const endTime = Date.now()

    expect(action).toHaveBeenCalledTimes(3)
    // Check if the total time taken is in the expected range considering concurrency limit
    expect(endTime - startTime).toBeGreaterThanOrEqual(200) // At least two batches of 100ms each
  })

  it('does not execute the same action for a given ID more than once', async () => {
    const action = jest.fn().mockResolvedValue(undefined)
    const concurrentAction = concurrentifyAction(action, 2)

    await Promise.all([
      concurrentAction('1'),
      concurrentAction('1') // This should not cause a second execution
    ])

    expect(action).toHaveBeenCalledTimes(1)
  })

  it('queues actions correctly when exceeding the concurrency limit', async () => {
    let activeCount = 0
    const action = jest.fn().mockImplementation(async (id) => {
      activeCount++
      expect(activeCount).toBeLessThanOrEqual(2) // Ensure no more than 2 active at a time
      await new Promise((resolve) => setTimeout(resolve, 50))
      activeCount--
    })
    const concurrentAction = concurrentifyAction(action, 2)

    await Promise.all([
      concurrentAction('1'),
      concurrentAction('2'),
      concurrentAction('3'),
      concurrentAction('4') // These should be queued
    ])

    expect(action).toHaveBeenCalledTimes(4)
  })
})
