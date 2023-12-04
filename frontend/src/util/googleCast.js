export function getCastSession() {
  if (!window.cast) {
    return undefined;
  }
  return window.cast.framework.CastContext.getInstance().getCurrentSession();
}

export function formatError(error) {
  const { ErrorCode } = window.chrome.cast;
  switch (error) {
    case ErrorCode.API_NOT_INITIALIZED:
      return `The API is not initialized. ${error.description || ""}`;
    case ErrorCode.CANCEL:
      return `The operation was canceled by the user. ${error.description || ""}`;
    case ErrorCode.CHANNEL_ERROR:
      return `A channel to the receiver is not available. ${error.description || ""}`;
    case ErrorCode.EXTENSION_MISSING:
      return `The Cast extension is not available. ${error.description || ""}`;
    case ErrorCode.INVALID_PARAMETER:
      return `The parameters to the operation were not valid. ${error.description || ""}`;
    case ErrorCode.RECEIVER_UNAVAILABLE:
      return `No receiver was compatible with the session request. ${error.description || ""}`;
    case ErrorCode.SESSION_ERROR:
      return `A session could not be created, or a session was invalid. ${error.description || ""}`;
    case ErrorCode.TIMEOUT:
      return `The operation timed out. ${error.description || ""}`;
    default:
      return error;
  }
}
