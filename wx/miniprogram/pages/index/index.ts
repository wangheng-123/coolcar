import { IAppOption } from "../../appoption";
import {routing} from "../../utils/routing";
// import {getUserInfo} from "../../utils/wxapi";

Page({
    isPageShowing:false,
    data: {
        avatarURL:'',
        settings:{
            skew: 0,
            rotate: 0,
            showLocation: true,
            showScale: true,
            subKey: '',
            layerStyle: -1,
            enableZoom: true,
            enableScroll: true,
            enableRotate: false,
            showCompass: false,
            enable3D: false,
            enableOverlooking: false,
            enableSatellite: false,
            enableTraffic: false,
        },
        location:{
            latitude:31,
            longitude:120,
        },
        scale:10,
        markers:[
            {
                iconPath:"/resources/car.png",
                id:0,
                latitude: 23.099994,
                longitude: 113.324520,
                width:50,
                height:50
            },
            {
                iconPath:"/resources/car.png",
                id:1,
                latitude: 23.099994,
                longitude: 114.324520,
                width:50,
                height:50
            },
        ],
    },
     async onLoad() {
         const userInfo = await getApp<IAppOption>().globalData.userInfo
         if (userInfo){
             console.log(userInfo.avatarUrl)
             this.setData({
                 avatarURL: userInfo.avatarUrl,
             })
         }
     },

    onMyLocationTap(){
        wx.getLocation({
            type:'gcj02',
            success:res=>{
                this.setData({
                    location: {
                        latitude: res.latitude,
                        longitude:res.longitude,
                    }
                })
            },
            fail:()=>{
                wx.showToast({
                    icon:'none',
                    title:'请前往设置页授权',
                })
            }
        })
    },

    onScanTap(){
        console.log("1")
      wx.scanCode({
          success:  ()=> {
              this.selectComponent('#licModal').showModal()
              //TODO get car id from scan result
              const carID = 'car123'
              // const redirectURL = `/pages/lock/lock?car_id=${carID}`
              const redirectURL = routing.lock({
                  car_id: carID
              })
              wx.navigateTo({
                  //url:`/pages/register/register?redirect=${encodeURIComponent(redirectURL)}`,
                  url: routing.register({
                      redirectURL: redirectURL,
                  })
              })
          },
          fail:console.error,
      })
        console.log("failde")
    },

    onShow(): void | Promise<void> {
        this.isPageShowing = true
    },

    onHide(): void | Promise<void> {
        this.isPageShowing=false
    },

    onMyTripsTap() {
        wx.navigateTo({
            // url:'/pages/mytrips/mytrips',
            url:routing.mytrips(),
        })
    },

    moveCars(){
        const map = wx.createMapContext("map");
        const dest={
            latitude: 23.099994,
            longitude: 113.324520,
        }
        const moveCar = ()=>{
            dest.latitude +=0.1
            dest.longitude +=0.1
            map.translateMarker({
                destination:{
                    latitude: dest.latitude+1,
                    longitude:dest.longitude+1,
                },
                markerId:0,
                autoRotate:false,
                rotate:0,
                duration:5000,
                animationEnd: ()=>{
                    if (this.isPageShowing){
                        moveCar()
                    }
                },
            })
        }
        moveCar()
    },
})