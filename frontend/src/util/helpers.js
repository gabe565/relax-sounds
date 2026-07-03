export const wait = (timeout) => new Promise((resolve) => setTimeout(resolve, timeout));

export const once = (fn) => {
  let promise;
  return () => (promise ??= fn());
};

export const toUrlSafeBase64 = (b64) => {
  let out = b64.replaceAll("+", "-").replaceAll("/", "_");
  while (out.endsWith("=")) out = out.slice(0, -1);
  return out;
};

export const fromUrlSafeBase64 = (str) => {
  str = str.replaceAll("-", "+").replaceAll("_", "/");
  return str + "==".slice(0, (4 - (str.length % 4)) % 4);
};

export const streamToArrayBuffer = async (stream) => {
  const blob = await new Response(stream).blob();
  return blob.arrayBuffer();
};

export const detectEncoding = (bytes) =>
  bytes[0] === 0x1f && bytes[1] === 0x8b ? "gzip" : "deflate";

export const compress = async (data, encoding = "deflate") => {
  const stream = new Blob([data]).stream().pipeThrough(new CompressionStream(encoding));
  const buffer = await streamToArrayBuffer(stream);
  const bytes = new Uint8Array(buffer);
  let binary = "";
  for (const b of bytes) binary += String.fromCodePoint(b);
  return toUrlSafeBase64(btoa(binary));
};

export const decompress = async (data, encoding = "") => {
  const binary = atob(fromUrlSafeBase64(data));
  const bytes = Uint8Array.from(binary, (c) => c.codePointAt(0));
  if (!encoding) {
    encoding = detectEncoding(bytes);
  }
  const stream = new Blob([bytes]).stream().pipeThrough(new DecompressionStream(encoding));
  return new Response(stream).text();
};
