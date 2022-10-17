module.exports = {
    apps: [
        {
            name: 'ICS Seafood Landing',
            script: 'npm',
            args: 'start',
            cwd: "/root/ics-seafood-landing",
            instances: 1,
            autorestart: true,
            watch: false,
            max_memory_restart: '1G',
            env: {
                NODE_ENV: "production",
            },
        },
        {
            name: 'ICS Seafood API',
            script: 'npm',
            args: 'start',
            cwd: "/root/ics-seafood-api",
            instances: 1,
            autorestart: true,
            watch: false,
            max_memory_restart: '1G',
            env: {
                NODE_ENV: "production",
            },
        },
    ]
};