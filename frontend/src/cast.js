import { ApiPath } from "@/config/api";
import { usePlayer } from "@/plugins/store/player";

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
    usePlayer().initializeCastApi(await castAppId);
  }
});
