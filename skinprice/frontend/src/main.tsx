import React from 'react'
import {createRoot} from 'react-dom/client'
import './styles.css'
import App from './App'
import { installGlobalErrorLogging } from './shared/lib/logging/logger'

const container = document.getElementById('root')

const root = createRoot(container!)

installGlobalErrorLogging()

root.render(
    <React.StrictMode>
        <App/>
    </React.StrictMode>
)
