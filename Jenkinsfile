@Library('libpipelines@master') _

hose {
    EMAIL = 'qa'
    MODULE = 'valkiria'
    REPOSITORY = 'valkiria'
    SLACKTEAM = 'stratiopaas'
    BUILDTOOL = 'make'
    DEVTIMEOUT = 10

    DEV = { config ->        
        doCompile(config)
        doUT(config)
        doPackage(conf: config, skipOnPR: true)
        doDeploy(config)
     }
}
