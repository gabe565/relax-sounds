/* eslint-disable no-unused-vars */
const initializeCastApi = () => {
  const { cast, chrome } = window;

  cast.framework.CastContext.getInstance().setOptions({
    receiverApplicationId: chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
    autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED,
  });

  const remotePlayer = new cast.framework.RemotePlayer();
  const remotePlayerController = new cast.framework.RemotePlayerController(remotePlayer);
  remotePlayerController.addEventListener(
    cast.framework.RemotePlayerEventType.IS_CONNECTED_CHANGED,
    (e) => console.log(e),
  );

  remotePlayerController.addEventListener(
    cast.framework.RemotePlayerEventType.ANY_CHANGE,
    (e) => console.log(e),
  );
};

// eslint-disable-next-line no-underscore-dangle
// eslint-disable-next-line dot-notation
window['__onGCastApiAvailable'] = (isAvailable) => {
  if (isAvailable) {
    initializeCastApi();
  }
};
