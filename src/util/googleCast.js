import { encode } from './shareUrl';
import soundConfig from '../data/sounds.json';

export function getCastSession() {
  return window.cast.framework.CastContext.getInstance().getCurrentSession();
}

export function castPlaylist(castSession, playlist) {
  const { sounds } = encode(playlist);
  const url = `${window.location.origin}/mix/${sounds}`;
  const mediaInfo = new window.chrome.cast.media.MediaInfo(url, 'music');

  mediaInfo.metadata = new window.chrome.cast.media.MusicTrackMediaMetadata();
  const playlistIds = playlist.sounds.map((s) => s.id);
  mediaInfo.metadata.title = soundConfig
    .filter((s) => playlistIds.includes(s.id))
    .map((s) => s.name)
    .sort((a, b) => a.localeCompare(b))
    .join(', ');
  mediaInfo.metadata.artist = 'Relax Sounds';

  const request = new window.chrome.cast.media.LoadRequest(mediaInfo);
  return castSession.loadMedia(request);
}
