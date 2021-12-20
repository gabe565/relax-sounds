import { encode } from './shareUrl';

export function getCastSession() {
  return window.cast.framework.CastContext.getInstance().getCurrentSession();
}

export function castPlaylist(castSession, playlist) {
  const { sounds } = encode(playlist);
  const url = `${window.location.origin}/mix/${sounds}`;
  const mediaInfo = new window.chrome.cast.media.MediaInfo(url, 'music');
  const request = new window.chrome.cast.media.LoadRequest(mediaInfo);
  return castSession.loadMedia(request);
}
