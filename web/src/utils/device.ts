/**
 * 获取设备信息
 */

// 解析User Agent获取设备信息
export function getDeviceInfo(): string {
  const ua = navigator.userAgent

  // 检测移动设备
  if (/Android/i.test(ua)) {
    // 尝试提取Android设备型号
    const match = ua.match(/Android.*?;\s*([^;)]+)/i)
    if (match && match[1]) {
      // 清理设备名称
      let deviceName = match[1].trim()
      // 移除 Build/ 之后的部分
      deviceName = deviceName.split(/Build|;/)[0].trim()
      return deviceName
    }
    return 'Android设备'
  }

  if (/iPhone/i.test(ua)) {
    // 提取iPhone型号
    const match = ua.match(/iPhone\s*(\w+)?/i)
    if (match) {
      return match[0]
    }
    return 'iPhone'
  }

  if (/iPad/i.test(ua)) {
    return 'iPad'
  }

  if (/Macintosh/i.test(ua)) {
    return 'Mac'
  }

  if (/Windows/i.test(ua)) {
    return 'Windows PC'
  }

  if (/Linux/i.test(ua)) {
    return 'Linux设备'
  }

  // 默认返回浏览器信息
  if (/Chrome/i.test(ua)) {
    return 'Chrome浏览器'
  }
  if (/Firefox/i.test(ua)) {
    return 'Firefox浏览器'
  }
  if (/Safari/i.test(ua)) {
    return 'Safari浏览器'
  }
  if (/Edge/i.test(ua)) {
    return 'Edge浏览器'
  }

  return '未知设备'
}

// 获取浏览器信息
export function getBrowserInfo(): string {
  const ua = navigator.userAgent

  if (/Edg/i.test(ua)) {
    return 'Edge'
  }
  if (/Chrome/i.test(ua)) {
    return 'Chrome'
  }
  if (/Firefox/i.test(ua)) {
    return 'Firefox'
  }
  if (/Safari/i.test(ua)) {
    return 'Safari'
  }

  return 'Unknown'
}

// 获取操作系统信息
export function getOSInfo(): string {
  const ua = navigator.userAgent

  if (/Android/i.test(ua)) {
    return 'Android'
  }
  if (/iPhone|iPad|iPod/i.test(ua)) {
    return 'iOS'
  }
  if (/Macintosh/i.test(ua)) {
    return 'macOS'
  }
  if (/Windows/i.test(ua)) {
    return 'Windows'
  }
  if (/Linux/i.test(ua)) {
    return 'Linux'
  }

  return 'Unknown'
}
