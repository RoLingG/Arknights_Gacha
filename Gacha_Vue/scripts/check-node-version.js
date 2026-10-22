const REQUIRED = '22.23.4'

const current = process.versions.node

if (current !== REQUIRED) {
  console.error(
    `\nNode.js 版本必须为 ${REQUIRED}，当前为 ${current}。\n` +
      `请切换到正确的版本后重试，例如：nvm install ${REQUIRED} && nvm use ${REQUIRED}\n`,
  )
  process.exit(1)
}