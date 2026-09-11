import { defineStore } from 'pinia'

export type NoticeKind = 'ok' | 'note' | 'fail'

type Notice = {
  id: number
  text: string
  kind: NoticeKind
}

const NOTICE_MS = 2500

let noticeTimer: ReturnType<typeof setTimeout> | null = null
let confirmWait: ((ok: boolean) => void) | null = null
let noticeSeq = 0

export const useFeedback = defineStore('feedback', {
  state: () => ({
    flash: null as Notice | null,
    confirmText: null as string | null,
  }),
  actions: {
    notice(text: string, kind: NoticeKind = 'note') {
      if (noticeTimer) {
        clearTimeout(noticeTimer)
        noticeTimer = null
      }
      this.flash = { id: ++noticeSeq, text, kind }
      noticeTimer = setTimeout(() => {
        this.flash = null
        noticeTimer = null
      }, NOTICE_MS)
    },
    dismissNotice() {
      if (noticeTimer) {
        clearTimeout(noticeTimer)
        noticeTimer = null
      }
      this.flash = null
    },
    confirm(text: string): Promise<boolean> {
      if (confirmWait) {
        confirmWait(false)
        confirmWait = null
      }
      this.confirmText = text
      return new Promise((resolve) => {
        confirmWait = resolve
      })
    },
    settle(ok: boolean) {
      this.confirmText = null
      const wait = confirmWait
      confirmWait = null
      wait?.(ok)
    },
  },
})
