export const wait = (timeout) => new Promise((resolve) => setTimeout(resolve, timeout));

export const toUrlSafeBase64 = (b64) =>
  b64.replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");

export const fromUrlSafeBase64 = (str) => {
  str = str.replace(/-/g, "+").replace(/_/g, "/");
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
  for (let i = 0; i < bytes.length; i++) binary += String.fromCharCode(bytes[i]);
  return toUrlSafeBase64(btoa(binary));
};

export const decompress = async (data, encoding = "") => {
  const binary = atob(fromUrlSafeBase64(data));
  const bytes = Uint8Array.from(binary, (c) => c.charCodeAt(0));
  if (!encoding) {
    encoding = detectEncoding(bytes);
  }
  const stream = new Blob([bytes]).stream().pipeThrough(new DecompressionStream(encoding));
  return new Response(stream).text();
};
