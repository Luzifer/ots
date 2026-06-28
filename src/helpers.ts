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

function durationToSeconds(duration) {
  const regex = /^(\d+)([smhd])$/
  const match = typeof duration === 'string' && duration.match(regex)
  if (!match) {
    return duration
  }

  const value = parseInt(match[1], 10)
  const unit = match[2]

  switch (unit) {
  case 's':
    return value
  case 'm':
    return value * 60
  case 'h':
    return value * 3600
  case 'd':
    return value * 86400
  }

  return duration // Fallback: return as-is
}

export {
  bytesToHuman,
  durationToSeconds,
}
