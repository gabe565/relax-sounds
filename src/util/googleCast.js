import { encode } from './shareUrl';

export function getCastSession() {
  return window.cast.framework.CastContext.getInstance().getCurrentSession();
}

export function castPlaylist(castSession, playlist) {
  const { sounds } = encode(playlist);
  const url = `${window.location.origin}/mix/${playlist.name}/${sounds}`;
  const mediaInfo = new window.chrome.cast.media.MediaInfo(url, 'music');

  mediaInfo.metadata = new window.chrome.cast.media.MusicTrackMediaMetadata();
  mediaInfo.metadata.title = playlist.name;
  mediaInfo.metadata.artist = 'Relax Sounds';

  const request = new window.chrome.cast.media.LoadRequest(mediaInfo);
  return castSession.loadMedia(request);
}
