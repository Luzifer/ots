import { appendFile } from 'node:fs/promises'

const allowedPermissions = new Set(['admin', 'maintain', 'write'])
const apiVersion = '2026-03-10'
const requiredChangedFile = 'i18n.yaml'

type PullRequest = {
  head: {
    ref: string
    repo: {
      full_name: string
    }
    sha: string
  }
  state: string
}

type PullRequestFile = {
  filename: string
}

type RepositoryPermission = {
  permission: string
}

async function api<T>(path: string): Promise<T> {
  const response = await fetch(`${requireEnv('GITHUB_API_URL')}${path}`, {
    headers: {
      Accept: 'application/vnd.github+json',
      Authorization: `Bearer ${requireEnv('GITHUB_TOKEN')}`,
      'X-GitHub-Api-Version': apiVersion,
    },
  })

  if (!response.ok) {
    throw new Error(`GitHub API request failed: ${path}: ${response.status} ${response.statusText}`)
  }

  return response.json() as Promise<T>
}

async function fetchChangedFiles(repository: string, pullRequestNumber: string): Promise<Set<string>> {
  const changedFiles = new Set<string>()

  for (let page = 1; ; page++) {
    const pageFiles = await api<PullRequestFile[]>(
      `/repos/${repository}/pulls/${pullRequestNumber}/files?per_page=100&page=${page}`,
    )

    if (pageFiles.length === 0) {
      break
    }

    for (const file of pageFiles) {
      changedFiles.add(file.filename)
    }
  }

  return changedFiles
}

async function main(): Promise<void> {
  const commentAuthor = requireEnv('COMMENT_AUTHOR')
  const outputFile = requireEnv('GITHUB_OUTPUT')
  const pullRequestNumber = requireEnv('PR_NUMBER')
  const repository = requireEnv('GITHUB_REPOSITORY')

  const { permission } = await api<RepositoryPermission>(
    `/repos/${repository}/collaborators/${commentAuthor}/permission`,
  )

  if (!allowedPermissions.has(permission)) {
    throw new Error(`User ${commentAuthor} has insufficient permission: ${permission}`)
  }

  const pullRequest = await api<PullRequest>(`/repos/${repository}/pulls/${pullRequestNumber}`)

  if (pullRequest.state !== 'open') {
    throw new Error(`Pull request #${pullRequestNumber} is not open: ${pullRequest.state}`)
  }

  const changedFiles = await fetchChangedFiles(repository, pullRequestNumber)

  if (!changedFiles.has(requiredChangedFile)) {
    throw new Error(`Pull request #${pullRequestNumber} does not change ${requiredChangedFile}`)
  }

  await writeOutput(outputFile, 'authorized', 'true')
  await writeOutput(outputFile, 'head_ref', pullRequest.head.ref)
  await writeOutput(outputFile, 'head_repo', pullRequest.head.repo.full_name)
  await writeOutput(outputFile, 'head_sha', pullRequest.head.sha)
}

function requireEnv(name: string): string {
  const value = process.env[name]

  if (!value) {
    throw new Error(`Missing required environment variable: ${name}`)
  }

  return value
}

async function writeOutput(outputFile: string, name: string, value: string): Promise<void> {
  await appendFile(outputFile, `${name}=${value}\n`)
}

main().catch((err: Error) => {
  console.error(err.message)
  process.exit(1)
})
