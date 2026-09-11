export type InstallSurface = 'hidden' | 'install' | 'ios-hint' | 'installed'

export type InstallContext = {
  standalone: boolean
  ios: boolean
  canPrompt: boolean
}

export function installSurface(ctx: InstallContext): InstallSurface {
  if (ctx.standalone) return 'installed'
  if (ctx.ios) return 'ios-hint'
  if (ctx.canPrompt) return 'install'
  return 'hidden'
}

export function isIosDevice(ua: string, maxTouchPoints: number): boolean {
  if (/iPhone|iPod|iPad/.test(ua)) return true
  return /Macintosh/.test(ua) && maxTouchPoints > 1
}
