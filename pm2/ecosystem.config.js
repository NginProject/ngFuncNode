module.exports = {
  apps : [{
    name: 'ngFuncNode',
    script: '../bin/ngFuncNode',

    // Options reference: https://pm2.io/doc/en/runtime/reference/ecosystem-file/
    args: '-a <YOUR_IP>',
    instances: 1,
    autorestart: true,
    watch: true,
    max_memory_restart: '1G',
  }]
};
