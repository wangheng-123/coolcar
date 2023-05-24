// // app.ts
// import { IAppOption } from "./appoption";
// import {getSetting, getUserInfo} from "./utils/wxapi";
// import {Coolcar} from "./service/request";
// // import {LoginResponse} from "./service/auth/v1";
//
// let resolveUserInfo: (value: WechatMiniprogram.UserInfo | PromiseLike<WechatMiniprogram.UserInfo>) => void
// let rejectUserInfo: (reason?: any) => void
//
// App<IAppOption>({
//     globalData: {
//         userInfo: new Promise((resolve, reject) => {
//             resolveUserInfo = resolve
//             rejectUserInfo = reject
//             console.log("globalData")
//         })
//     },
//     async onLaunch() {
//
//         // wx.request({
//         //     url:'http://localhost:8080/trip/trip123',
//         //     method:'GET',
//         //     success:console.log,
//         //     fail:console.error,
//         // })
//
//         // // 展示本地存储能力
//         // const logs = wx.getStorageSync('logs') || []
//         // logs.unshift(Date.now())
//         // wx.setStorageSync('logs', logs)
//
//         // 登录
//         await Coolcar.login()
//         // wx.login({
//         //     success: res => {
//         //         console.log(res.code)
//         //         // 发送 res.code 到后台换取 openId, sessionKey, unionId
//         //         wx.request({
//         //             url:'http://localhost:8080/v1/auth/login',
//         //             method: 'POST',
//         //             data:{
//         //                 code : res.code
//         //             },
//         //             success:res=>{
//         //                 console.log(res)
//         //                 console.log(res.data.accessToken)
//         //                 wx.request({
//         //                     url:'http://localhost:8080/v1/trip',
//         //                     method:'POST',
//         //                     data:{
//         //                         start:'abc',
//         //                     },
//         //                     header:{
//         //                         authorization:'Bearer '+res.data.accessToken
//         //                     }
//         //                 })
//         //             },
//         //             fail:console.error,
//         //         })
//         //     },
//         // })
//
//         //获取用户信息
//         try {
//             const setting = await getSetting()
//             if (setting.authSetting['scope.userInfo']){
//                 const userInfoRes = await getUserInfo()
//                 resolveUserInfo(userInfoRes.userInfo)
//             }
//         } catch (err){
//             rejectUserInfo(err)
//         }
//     },
//
//     resolveUserInfo(userInfo:WechatMiniprogram.UserInfo){
//         resolveUserInfo(userInfo)
//     }
// })





import { IAppOption } from "./appoption"
import { Coolcar } from "./service/request"
import { getSetting, getUserInfo } from "./utils/wxapi"

let resolveUserInfo: (value: WechatMiniprogram.UserInfo | PromiseLike<WechatMiniprogram.UserInfo>) => void
let rejectUserInfo: (reason?: any) => void

// app.ts
App<IAppOption>({
    globalData: {
        userInfo: new Promise((resolve, reject) => {
            resolveUserInfo = resolve
            rejectUserInfo = reject
        })
    },
    async onLaunch() {
        // 登录
        Coolcar.login()

        // 获取用户信息
        try {
            const setting = await getSetting()
            if (setting.authSetting['scope.userInfo']) {
                const userInfoRes = await getUserInfo()
                resolveUserInfo(userInfoRes.userInfo)
            }
        } catch (err) {
            rejectUserInfo(err)
        }
    },
    resolveUserInfo(userInfo: WechatMiniprogram.UserInfo) {
        resolveUserInfo(userInfo)
    }
})
