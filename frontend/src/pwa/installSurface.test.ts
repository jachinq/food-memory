import { describe, expect, it } from 'vitest'
import { installSurface, isIosDevice } from './installSurface'

describe('installSurface', () => {
  it('hides the install control when already running as 安装的食忆', () => {
    expect(
      installSurface({ standalone: true, ios: false, canPrompt: true }),
    ).toBe('installed')
  })

  it('shows iOS add-to-home-screen copy instead of a fake install button', () => {
    expect(
      installSurface({ standalone: false, ios: true, canPrompt: false }),
    ).toBe('ios-hint')
  })

  it('shows install when the browser can prompt', () => {
    expect(
      installSurface({ standalone: false, ios: false, canPrompt: true }),
    ).toBe('install')
  })

  it('hides install when the browser cannot prompt', () => {
    expect(
      installSurface({ standalone: false, ios: false, canPrompt: false }),
    ).toBe('hidden')
  })
})

describe('isIosDevice', () => {
  it('treats iPhone Safari as iOS', () => {
    expect(
      isIosDevice(
        'Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1',
        5,
      ),
    ).toBe(true)
  })

  it('treats iPadOS desktop UA with touch as iOS', () => {
    expect(
      isIosDevice(
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15',
        5,
      ),
    ).toBe(true)
  })

  it('does not treat Windows Chrome as iOS', () => {
    expect(
      isIosDevice(
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
        0,
      ),
    ).toBe(false)
  })
})
