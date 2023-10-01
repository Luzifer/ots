/**
 * Converts number of bytes into human format (524288 -> "512.0 KiB")
 * @param {Number} bytes Byte amount to convert into human readable format
 * @returns String
 */
function bytesToHuman(bytes) {
  for (const t of [
    { thresh: 1024 * 1024, unit: 'MiB' },
    { thresh: 1024, unit: 'KiB' },
  ]) {
    if (bytes > t.thresh) {
      return `${(bytes / t.thresh).toFixed(1)} ${t.unit}`
    }
  }

  return `${bytes} B`
}

export {
  bytesToHuman,
}
