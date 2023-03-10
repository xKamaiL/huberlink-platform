'use client'
import AuthService from '@/services/AuthService'
import { useEffect, useState } from 'react'
import { IUser } from '../types'
import { createProvider } from './index'

export type IUserActions =
  | {
      type: 'fetch-user'
    }
  | { type: 'open-dialog' }
  | { type: 'close-dialog' }
  | { type: 'logout' }

type IStatus = 'idle' | 'loading' | 'success'

export const [UserContextProvider, useUserDispatch, useUserSelector] =
  createProvider(() => {
    const [userData, setUserData] = useState<IUser | null>(null)
    const [status, setStatus] = useState<IStatus>('idle')

    useEffect(() => {
      if (status === 'idle') {
        dispatch({ type: 'fetch-user' })
      }
    }, [status])

    const dispatch = async (action: IUserActions) => {
      switch (action.type) {
        case 'logout':
          setUserData(null)
          return

        case 'fetch-user':
          setStatus('loading')
          await AuthService.me()
            .then((r) => {
              if (!r.success) {
                setUserData(null)
                return
              }
              setUserData(r.data)
            })
            .catch((err) => {
              setUserData(null)
            })
            .finally(() => {
              setStatus('success')
            })
          return
      }
    }
    const isLoggedIn = !!userData
    const state = {
      profile: userData,
      isLoggedIn,
      status,
    }
    return [state, dispatch]
  })
