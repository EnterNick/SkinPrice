import React from 'react'
import {createRoot} from 'react-dom/client'
import './styles.css'
import App from './App'
import { installGlobalErrorLogging } from './shared/lib/logging/logger'
import { applyTypographySettings } from './shared/lib/settings/appTypography'

const container = document.getElementById('root')

const root = createRoot(container!)

installGlobalErrorLogging()
applyTypographySettings()

root.render(
    <React.StrictMode>
        <App/>
    </React.StrictMode>
)
