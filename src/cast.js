const initializeCastApi = () => {
  const { cast, chrome } = window;
  console.log('SETUP');

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

// const castSdk = document.createElement('script');
// castSdk.async = true;
// castSdk.src = 'https://www.gstatic.com/cv/js/sender/v1/cast_sender.js?loadCastFramework=1';
// document.body.appendChild(castSdk);
