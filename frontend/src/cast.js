import { ApiPath } from "@/config/api";
import { usePlayer } from "@/plugins/store/player";
import { wait } from "@/util/helpers";

const fetchCastAppId = async () => {
  try {
    const resp = await fetch(ApiPath("/api/config"), { signal: AbortSignal.timeout(5000) });
    if (resp.ok) {
      const config = await resp.json();
      return config.castAppId;
    }
  } catch (error) {
    console.error("Failed to load config:", error);
  }
  return undefined;
};

const castAppId = fetchCastAppId();

globalThis.castApiAvailable.then(async (isAvailable) => {
  if (isAvailable) {
    // Workaround for __onGCastApiAvailable called before globalThis.cast is set
    let waitMs = 100;
    while (!globalThis.cast) {
      console.warn(`Cast is undefined. Retrying setup in ${waitMs}ms.`);
      await wait(waitMs);
      waitMs *= 2;
    }

    usePlayer().initializeCastApi(await castAppId);
  }
});
