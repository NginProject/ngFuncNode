module.exports = {
  apps : [{
    name: 'ngFuncNode',
    script: '../bin/ngFuncNode',

    args: '-a <YOUR_IP>',
    instances: 1,
    autorestart: true,
    watch: true,
    max_memory_restart: '1G',
  }]
};
