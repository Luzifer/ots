/**
 * Converts number of bytes into human format (524288 -> "512.0 KiB")
 * @param {number} bytes Byte amount to convert into human readable format
 * @returns string
 */
function bytesToHuman(bytes: number): string {
  for (const t of [
    { thresh: 1024 * 1024, unit: 'MiB' },
    { thresh: 1024, unit: 'KiB' },
  ]) {
    if (bytes > t.thresh) {
      return `${parseFloat((bytes / t.thresh).toFixed(1))} ${t.unit}`
    }
  }

  return `${bytes} B`
}

export {
  bytesToHuman,
}
