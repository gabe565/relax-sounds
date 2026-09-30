const { cast } = globalThis;

const TAG = "RelaxSounds";
const MIN_RETRY_MS = 2000;
const MAX_RETRY_MS = 30000;
const STALL_MS = 30000;
const HEALTHY_MS = 60000;

const { EventType, EndedReason } = cast.framework.events;
const { MessageType, RepeatMode } = cast.framework.messages;

const context = cast.framework.CastReceiverContext.getInstance();
const playerManager = context.getPlayerManager();

const params = new URLSearchParams(location.search);
const logger = cast.debug.CastDebugLogger.getInstance();
logger.setEnabled(params.has("debug"));
logger.loggerLevelByTags = { [TAG]: cast.framework.LoggerLevel.DEBUG };

let lastLoad = null;
let retryDelay = MIN_RETRY_MS;
let retryTimer = null;
let stallTimer = null;
let healthyTimer = null;

const clearTimers = () => {
  clearTimeout(retryTimer);
  clearTimeout(stallTimer);
  clearTimeout(healthyTimer);
  retryTimer = stallTimer = healthyTimer = null;
};

const scheduleReload = (reason) => {
  if (!lastLoad || retryTimer) return;
  clearTimeout(stallTimer);
  clearTimeout(healthyTimer);
  logger.warn(TAG, `Reloading in ${retryDelay}ms: ${reason}`);
  retryTimer = setTimeout(() => {
    retryTimer = null;
    lastLoad.queueData?.items?.forEach((item) => delete item.itemId);
    playerManager.load(lastLoad).catch((err) => {
      logger.error(TAG, `Reload failed: ${JSON.stringify(err)}`);
      scheduleReload("reload failed");
    });
  }, retryDelay);
  retryDelay = Math.min(retryDelay * 2, MAX_RETRY_MS);
};

playerManager.setMessageInterceptor(MessageType.LOAD, (request) => {
  clearTimers();
  if (request.queueData) {
    request.queueData.repeatMode = RepeatMode.REPEAT_OFF;
  }
  lastLoad = request;
  return request;
});

playerManager.setMessageInterceptor(MessageType.STOP, (request) => {
  clearTimers();
  lastLoad = null;
  return request;
});

playerManager.addEventListener(EventType.PLAYING, () => {
  clearTimeout(stallTimer);
  stallTimer = null;
  clearTimeout(healthyTimer);
  healthyTimer = setTimeout(() => {
    retryDelay = MIN_RETRY_MS;
  }, HEALTHY_MS);
});

playerManager.addEventListener(EventType.BUFFERING, ({ isBuffering }) => {
  clearTimeout(stallTimer);
  stallTimer = null;
  if (isBuffering && lastLoad) {
    stallTimer = setTimeout(() => scheduleReload("stalled"), STALL_MS);
  }
});

playerManager.addEventListener(EventType.ERROR, (event) => {
  logger.error(TAG, `Player error: ${event.detailedErrorCode} ${event.reason ?? ""}`);
  scheduleReload(`error ${event.detailedErrorCode}`);
});

playerManager.addEventListener(EventType.MEDIA_FINISHED, ({ endedReason }) => {
  logger.info(TAG, `Media finished: ${endedReason}`);
  if (endedReason === EndedReason.ERROR || endedReason === EndedReason.END_OF_STREAM) {
    scheduleReload(`ended ${endedReason}`);
  } else if (endedReason === EndedReason.STOPPED) {
    clearTimers();
    lastLoad = null;
  }
});

const options = new cast.framework.CastReceiverOptions();
options.supportedCommands = cast.framework.messages.Command.PAUSE;
context.start(options);
