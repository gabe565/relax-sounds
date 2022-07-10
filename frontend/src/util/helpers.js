// eslint-disable-next-line no-promise-executor-return
export const wait = (timeout) => new Promise((resolve) => setTimeout(resolve, timeout));
