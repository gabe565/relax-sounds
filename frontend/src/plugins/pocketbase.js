import PocketBase from "pocketbase";
import { ApiPath } from "@/config/api";

export const pb = new PocketBase(ApiPath());

/**
 * Parses a PocketBase error and returns a user-friendly message.
 * @param {Object} error - The error object from the PocketBase client.
 * @returns {string} - A formatted error message.
 */
export const getErrorMessage = (error) => {
  const response = error.response;
  if (response?.data && typeof response.data === "object") {
    const messages = [];
    for (const key in response.data) {
      if (response.data[key]?.message) {
        messages.push(`${key}: ${response.data[key].message}`);
      }
    }
    if (messages.length > 0) {
      return messages.join("\n");
    }
  }
  return response?.message || error.message || "An unexpected error occurred";
};
