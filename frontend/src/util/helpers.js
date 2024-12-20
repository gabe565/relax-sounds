import base64 from "base64-url";

export const wait = (timeout) => new Promise((resolve) => setTimeout(resolve, timeout));

export const compress = async (data, encoding = "gzip") => {
  const stream = new Blob([data]).stream();
  const compressedReadableStream = stream.pipeThrough(new CompressionStream(encoding));
  const compressedResponse = new Response(compressedReadableStream);
  const blob = await compressedResponse.blob();
  const buffer = await blob.arrayBuffer();
  const compressedBase64 = base64.escape(btoa(String.fromCharCode(...new Uint8Array(buffer))));
  return compressedBase64;
};

export const decompress = async (data, encoding = "gzip") => {
  data = atob(base64.unescape(data));
  data = Uint8Array.from(data, (m) => m.codePointAt(0));
  const stream = new Blob([data]).stream();
  const compressedReadableStream = stream.pipeThrough(new DecompressionStream(encoding));
  const resp = new Response(compressedReadableStream);
  const blob = await resp.blob();
  return await blob.text();
};
