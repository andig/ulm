(function(t){function e(e){for(var s,n,o=e[0],l=e[1],c=e[2],u=0,m=[];u<o.length;u++)n=o[u],Object.prototype.hasOwnProperty.call(i,n)&&i[n]&&m.push(i[n][0]),i[n]=0;for(s in l)Object.prototype.hasOwnProperty.call(l,s)&&(t[s]=l[s]);d&&d(e);while(m.length)m.shift()();return r.push.apply(r,c||[]),a()}function a(){for(var t,e=0;e<r.length;e++){for(var a=r[e],s=!0,o=1;o<a.length;o++){var l=a[o];0!==i[l]&&(s=!1)}s&&(r.splice(e--,1),t=n(n.s=a[0]))}return t}var s={},i={index:0},r=[];function n(e){if(s[e])return s[e].exports;var a=s[e]={i:e,l:!1,exports:{}};return t[e].call(a.exports,a,a.exports,n),a.l=!0,a.exports}n.m=t,n.c=s,n.d=function(t,e,a){n.o(t,e)||Object.defineProperty(t,e,{enumerable:!0,get:a})},n.r=function(t){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})},n.t=function(t,e){if(1&e&&(t=n(t)),8&e)return t;if(4&e&&"object"===typeof t&&t&&t.__esModule)return t;var a=Object.create(null);if(n.r(a),Object.defineProperty(a,"default",{enumerable:!0,value:t}),2&e&&"string"!=typeof t)for(var s in t)n.d(a,s,function(e){return t[e]}.bind(null,s));return a},n.n=function(t){var e=t&&t.__esModule?function(){return t["default"]}:function(){return t};return n.d(e,"a",e),e},n.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},n.p="";var o=window["webpackJsonp"]=window["webpackJsonp"]||[],l=o.push.bind(o);o.push=e,o=o.slice();for(var c=0;c<o.length;c++)e(o[c]);var d=l;r.push([0,"chunk-vendors"]),a()})({0:function(t,e,a){t.exports=a("9e0e")},2110:function(t,e,a){},3899:function(t,e,a){},"38d7":function(t,e,a){},"3c05":function(t,e,a){"use strict";a("3899")},"3ec5":function(t,e,a){"use strict";a("38d7")},"5aff":function(t,e,a){"use strict";a("8817")},8817:function(t,e,a){},"8aef":function(t,e,a){"use strict";a("c0cc")},"9e0e":function(t,e,a){"use strict";a.r(e);a("b19f"),a("ab8b"),a("7b17");var s=a("2b0e"),i=a("bc3a"),r=a.n(i),n=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"app d-flex flex-column justify-content-between overflow-hidden"},[a("div",{staticClass:"flex-grow-1 d-flex flex-column justify-content-between"},[a("nav",{staticClass:"navbar navbar-expand-lg navbar-dark bg-dark"},[a("div",{staticClass:"container"},[a("a",{staticClass:"navbar-brand",attrs:{href:"https://github.com/andig/evcc"}},[a("Logo",{staticClass:"logo"})],1),a("div",{staticClass:"d-flex"},[a("div",{staticClass:"d-flex"},[a("Notifications",{attrs:{notifications:t.notifications}}),t._m(0)],1),t._m(1)])])]),a("router-view",{staticClass:"flex-grow-1 d-flex flex-column justify-content-stretch"})],1),a("Footer",{attrs:{version:t.version,supporter:t.supporter}})],1)},o=[function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("button",{staticClass:"navbar-toggler",attrs:{type:"button","data-bs-toggle":"collapse","data-bs-target":"#navbarNavAltMarkup","aria-controls":"navbarNavAltMarkup","aria-expanded":"false","aria-label":"Toggle navigation"}},[a("span",{staticClass:"navbar-toggler-icon"})])},function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"collapse navbar-collapse flex-lg-grow-0",attrs:{id:"navbarNavAltMarkup"}},[a("ul",{staticClass:"navbar-nav"},[a("li",{staticClass:"nav-item"},[a("a",{staticClass:"nav-link",attrs:{href:"https://github.com/andig/evcc/discussions",target:"_blank"}},[t._v(" Support ")])])])])}],l=a("ad3d"),c=a("ecee"),d=a("a206"),u=a("9e52"),m=a("8475"),p=a("39f0"),v=a("8718"),h=a("7c8d"),f=a("c6b3"),g=a("8668"),b=a("a14b"),C=a("5dae"),_=a("8560"),y=a("bf13"),w=a("af2b"),S=a("f303"),x=a("7116"),k=a("fdca"),N=a("91fb"),D=a("184c"),M=a("ba01"),P=a("6c06"),T=a("adbc");c["c"].add(d["faArrowDown"],u["faArrowUp"],m["faBatteryEmpty"],p["faBatteryFull"],v["faBatteryHalf"],h["faBatteryQuarter"],f["faBatteryThreeQuarters"],g["faChevronDown"],b["faChevronUp"],C["faClock"],_["faExclamationTriangle"],y["faLeaf"],w["faSun"],S["faTemperatureHigh"],x["faTemperatureLow"],k["faThermometerHalf"],N["faHeart"],D["faGift"],M["faBox"],T["faExclamationCircle"],P["faStar"]),s["a"].component("fa-icon",l["a"]);var B=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("svg",{attrs:{viewBox:"0 0 122 35",xmlns:"http://www.w3.org/2000/svg","fill-rule":"evenodd","clip-rule":"evenodd","stroke-linejoin":"round","stroke-miterlimit":"2"}},[a("path",{attrs:{d:"M13.082 29.071a12.384 12.384 0 01-9-3.42 12.192 12.192 0 01-3.54-9.12v-.64a15.394 15.394 0 011.47-6.83 10.825 10.825 0 014.17-4.64 11.64 11.64 0 016.15-1.63 10.45 10.45 0 018.21 3.26c2 2.194 3 5.297 3 9.31v2.76H7.382a6.348 6.348 0 002 4 5.997 5.997 0 004.16 1.49 7.305 7.305 0 006.1-2.84l3.31 3.73a10 10 0 01-4.13 3.39 13.309 13.309 0 01-5.74 1.18zm-.77-20.84a4.216 4.216 0 00-3.26 1.37 7.141 7.141 0 00-1.6 3.91h9.39v-.55a5.005 5.005 0 00-1.22-3.49 4.304 4.304 0 00-3.31-1.24zM36.452 20.331l4.7-17.09h7l-8.48 25.36h-6.44l-8.52-25.36h7l4.74 17.09zM85.542 23.611a4.444 4.444 0 003-1 3.638 3.638 0 001.22-2.75h6.32a8.668 8.668 0 01-1.4 4.73 9.145 9.145 0 01-3.79 3.3 11.736 11.736 0 01-5.29 1.19 10.912 10.912 0 01-8.54-3.46c-2.087-2.3-3.13-5.483-3.13-9.55v-.45c0-3.9 1.033-7.016 3.1-9.35a10.868 10.868 0 018.51-3.5c2.791-.134 5.524.84 7.6 2.71a9.626 9.626 0 012.9 7.21h-6.3a4.663 4.663 0 00-1.2-3.22 4.005 4.005 0 00-3.08-1.24 4.068 4.068 0 00-3.56 1.73c-.8 1.15-1.2 3-1.2 5.6v.7c0 2.61.39 4.49 1.19 5.63a4.092 4.092 0 003.65 1.72zM110.422 23.611a4.454 4.454 0 003-1 3.63 3.63 0 001.21-2.75h6.33a8.668 8.668 0 01-1.4 4.73 9.143 9.143 0 01-3.73 3.3 11.76 11.76 0 01-5.29 1.18 10.912 10.912 0 01-8.54-3.46c-2.087-2.3-3.13-5.483-3.13-9.55v-.45c0-3.9 1.033-7.016 3.1-9.35a10.85 10.85 0 018.57-3.49 10.575 10.575 0 017.6 2.71 9.598 9.598 0 012.91 7.21h-6.33a4.651 4.651 0 00-1.21-3.22 4.492 4.492 0 00-6.64.49c-.8 1.15-1.21 3-1.21 5.6v.7c0 2.607.4 4.484 1.2 5.63a4.09 4.09 0 003.56 1.72z",fill:"#fff","fill-rule":"nonzero"}}),a("path",{attrs:{d:"M58.462.751h9.22l-6.14 12.3h6.15l-11.53 21.51 2.3-15.36h-7.68l7.68-18.45z",fill:"#0fdd42","fill-rule":"nonzero"}}),a("path",{attrs:{fill:"none",d:"M-24.458-22.109h170v76h-170z"}})])},E=[],V={name:"Logo"},j=V,A=a("2877"),$=Object(A["a"])(j,B,E,!1,null,null,null),z=$.exports,L=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("footer",{staticClass:"container"},[a("div",{staticClass:"py-3 py-md-5 mt-3 mt-md-5 border-top"},[a("div",{staticClass:"d-flex justify-content-between"},[a("Version",t._b({},"Version",t.version,!1)),a("Supporter",{attrs:{supporter:t.supporter}})],1)])])},O=[],U=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[t.newVersionAvailable?a("button",{staticClass:"btn btn-link ps-0 text-decoration-none link-dark text-nowrap",attrs:{href:"#","data-bs-toggle":"modal","data-bs-target":"#updateModal"}},[a("fa-icon",{staticClass:"icon me-1",attrs:{icon:"gift"}}),t._v(" Update"),a("span",{staticClass:"d-none d-sm-inline"},[t._v(" verfügbar")]),t._v(": "+t._s(t.available)+" ")],1):a("a",{staticClass:"btn btn-link ps-0 text-decoration-none link-dark text-nowrap",attrs:{href:t.releaseNotesUrl(t.installed),target:"_blank"}},[t._v(" Version "+t._s(t.installed)+" ")]),a("div",{staticClass:"modal fade",attrs:{id:"updateModal",tabindex:"-1",role:"dialog","aria-hidden":"true"}},[a("div",{staticClass:"modal-dialog modal-dialog-centered modal-dialog-scrollable",attrs:{role:"document"}},[a("div",{staticClass:"modal-content"},[t._m(0),a("div",{staticClass:"modal-body"},[t.updateStarted?a("div",[a("p",[t._v("Nach der Aktualisierung wird evcc neu gestartet.")]),a("div",{staticClass:"progress my-3"},[a("div",{staticClass:"progress-bar progress-bar-striped progress-bar-animated",style:{width:t.uploadProgress+"%"},attrs:{role:"progressbar"}})]),a("p",[t._v(t._s(t.updateStatus)+t._s(t.uploadMessage))])]):a("div",[a("p",[a("small",[t._v("Aktuell installierte Version: "+t._s(t.installed))])]),t.releaseNotes?a("div",{domProps:{innerHTML:t._s(t.releaseNotes)}}):a("p",[t._v(" Keine Releasenotes verfügbar. Mehr Informationen zur neuen Version findest du "),a("a",{attrs:{href:t.releaseNotesUrl(t.available)}},[t._v("hier")]),t._v(". ")])])]),a("div",{staticClass:"modal-footer d-flex justify-content-between"},[a("button",{staticClass:"btn btn-outline-secondary",attrs:{type:"button",disabled:t.updateStarted,"data-bs-dismiss":"modal"}},[t._v(" Abbrechen ")]),a("div",[t.hasUpdater?a("button",{staticClass:"btn btn-primary",attrs:{type:"button",disabled:t.updateStarted},on:{click:t.update}},[t.updateStarted?a("span",[a("span",{staticClass:"spinner-border spinner-border-sm",attrs:{role:"status","aria-hidden":"true"}}),t._v(" Akualisieren ")]):a("span",[t._v("Jetzt aktualisieren")])]):a("a",{staticClass:"btn btn-primary",attrs:{href:t.releaseNotesUrl(t.available)}},[t._v(" Download ")])])])])])])])},I=[function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"modal-header"},[a("h5",{staticClass:"modal-title"},[t._v("Update verfügbar")]),a("button",{staticClass:"btn-close",attrs:{type:"button","data-bs-dismiss":"modal","aria-label":"Close"}})])}],W={name:"Version",props:{installed:String,available:String,releaseNotes:String,hasUpdater:Boolean,uploadMessage:String,uploadProgress:Number,supporter:Boolean},data:function(){return{updateStarted:!1,updateStatus:""}},methods:{update:async function(){try{await r.a.post("update"),this.updateStatus="Aktualisierung gestartet: ",this.updateStarted=!0}catch(t){this.updateStatus="Aktualisierung nicht möglich: "+t}},releaseNotesUrl:function(t){return"https://github.com/andig/evcc/releases/tag/"+t}},computed:{newVersionAvailable:function(){return this.available&&"[[.Version]]"!=this.installed&&"0.0.1-alpha"!=this.installed&&this.available!=this.installed}}},H=W,F=(a("3ec5"),Object(A["a"])(H,U,I,!1,null,"5f5c4634",null)),K=F.exports,R=function(){var t=this,e=t.$createElement,a=t._self._c||e;return t.supporter?a("div",{ref:"supporter",staticClass:"btn btn-link pe-0 text-decoration-none link-dark text-nowrap supporter-button",on:{click:function(e){return e.stopPropagation(),e.preventDefault(),t.surprise(e)}}},[a("span",[t._v("Supporter")]),a("fa-icon",{staticClass:"icon ms-1",attrs:{icon:"star"}})],1):a("a",{staticClass:"btn btn-link pe-0 text-decoration-none link-dark text-nowrap",attrs:{href:"https://github.com/sponsors/andig",target:"_blank"}},[a("span",{staticClass:"d-none d-sm-inline"},[t._v("Projekt ")]),t._v("unterstützen "),a("fa-icon",{staticClass:"icon ms-1",attrs:{icon:"heart"}})],1)},Z=[],J=a("7129"),G={name:"Supporter",props:{supporter:Boolean},methods:{surprise:function(){console.log(this.$refs.supporter);const{top:t,height:e,left:a,width:s}=this.$refs.supporter.getBoundingClientRect(),i=(a+s/2)/window.innerWidth,r=(t+e/2)/window.innerHeight,n={x:i,y:r};Object(J["a"])({origin:n,angle:90+35*Math.random(),particleCount:75+50*Math.random(),spread:50+50*Math.random(),drift:-.5,scalar:1.3,colors:["#0d6efd","#0fdd42","#408458","#4923BA","#5BC8EC","#C54482","#CC444A","#EE8437","#F7C144","#FFFD54"]})}}},q=G,Q=(a("5aff"),Object(A["a"])(q,R,Z,!1,null,"c13137dc",null)),Y=Q.exports,X={name:"Footer",components:{Version:K,Supporter:Y},props:{version:Object,supporter:Boolean}},tt=X,et=Object(A["a"])(tt,L,O,!1,null,null,null),at=et.exports,st=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("button",{directives:[{name:"show",rawName:"v-show",value:t.iconVisible,expression:"iconVisible"}],staticClass:"btn btn-link text-decoration-none link-light text-nowrap",attrs:{href:"#","data-bs-toggle":"modal","data-bs-target":"#notificationModal"}},[a("fa-icon",{class:t.iconClass,attrs:{icon:"exclamation-triangle"}})],1),a("div",{staticClass:"modal fade",attrs:{id:"notificationModal",tabindex:"-1",role:"dialog","aria-hidden":"true"}},[a("div",{staticClass:"modal-dialog modal-dialog-centered modal-dialog-scrollable",attrs:{role:"document"}},[a("div",{staticClass:"modal-content"},[t._m(0),a("div",{staticClass:"modal-body"},t._l(t.notifications,(function(e,s){return a("p",{key:s,staticClass:"d-flex align-items-baseline"},[a("fa-icon",{staticClass:"flex-grow-0 d-block",class:{"text-danger":"error"===e.type,"text-warning":"warn"===e.type},attrs:{icon:"exclamation-triangle"}}),a("span",{staticClass:"flex-grow-1 px-2 py-1"},[t._v(t._s(e.message))])],1)})),0),a("div",{staticClass:"modal-footer"},[a("button",{staticClass:"btn btn-outline-secondary",attrs:{type:"button","data-bs-dismiss":"modal","aria-label":"Close"},on:{click:t.clear}},[t._v(" Meldungen entfernen ")])])])])])])},it=[function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"modal-header"},[a("h5",{staticClass:"modal-title"},[t._v("Meldungen")]),a("button",{staticClass:"btn-close",attrs:{type:"button","data-bs-dismiss":"modal","aria-label":"Close"}})])}],rt={data:function(){return{fmtLimit:100,fmtDigits:1}},methods:{round:function(t,e){var a=10**e;return(Math.round(t*a)/a).toFixed(e)},fmt:function(t){return void 0===t||null===t?0:(t=Math.abs(t),t>=this.fmtLimit?this.round(t/1e3,this.fmtDigits):this.round(t,0))},fmtUnit:function(t){return Math.abs(t)>=this.fmtLimit?"k":""},fmtDuration:function(t){if(t<=0||null==t)return"—";var e="0"+t%60,a="0"+Math.floor(t/60)%60,s=""+Math.floor(t/3600);return s.length<2&&(s="0"+s),s+":"+a.substr(-2)+":"+e.substr(-2)},fmtShortDuration:function(t){if(t<=0||null==t)return"—";var e,a=Math.floor(t/60)%60,s=Math.floor(t/3600);if(s>=1)a="0"+a,e=s+":"+a.substr(-2);else{var i="0"+t%60;e=a+":"+i.substr(-2)}return e},fmtShortDurationUnit:function(t){if(t<=0||null==t)return"";var e=Math.floor(t/3600);return e>=1?"h":"m"},fmtDayString:function(t){const e=""+t.getFullYear(),a=(""+(t.getMonth()+1)).padStart(2,"0"),s=(""+t.getDate()).padStart(2,"0");return`${e}-${a}-${s}`},fmtTimeString:function(t){const e=(""+t.getHours()).padStart(2,"0"),a=(""+t.getMinutes()).padStart(2,"0");return`${e}:${a}`},fmtAbsoluteDate:function(t){return new Intl.DateTimeFormat("default",{weekday:"short",hour:"numeric",minute:"numeric"}).format(t)}}},nt={name:"Notifications",props:{notifications:Array},computed:{iconVisible:function(){return this.notifications.length>0},iconClass:function(){return this.notifications.find(t=>"error"===t.type)?"text-danger":"text-warning"}},methods:{clear:function(){window.app.clear()}},mixins:[rt]},ot=nt,lt=Object(A["a"])(ot,st,it,!1,null,null,null),ct=lt.exports;function dt(t,e,a){const i=e.shift();t[i]||s["a"].set(t,i,{}),e.length?dt(t[i],e,a):a&&"object"===typeof a&&!Array.isArray(a)?t[i]={...t[i],...a}:t[i]=a}const ut={state:{loadpoints:[]},update:function(t){Object.keys(t).forEach((function(e){"function"===typeof window.app[e]?window.app[e]({message:t[e]}):dt(ut.state,e.split("."),t[e])}))}};var mt=ut,pt={name:"App",components:{Logo:z,Footer:at,Notifications:ct},data:function(){return{compact:!1,store:this.$root.$data.store,installedVersion:window.evcc.version}},methods:{connect:function(){const t=window.location,e="https:"==t.protocol?"wss:":"ws:",a=e+"//"+t.hostname+(t.port?":"+t.port:"")+t.pathname+"ws",s=new WebSocket(a),i=this;s.onerror=function(){s.close()},s.onclose=function(){window.setTimeout(i.connect,1e3)},s.onmessage=function(t){try{var e=JSON.parse(t.data);mt.update(e)}catch(a){window.app.error(a,t.data)}}}},computed:{version:function(){return{installed:this.installedVersion,available:this.store.state.availableVersion,releaseNotes:this.store.state.releaseNotes,hasUpdater:this.store.state.hasUpdater,uploadMessage:this.store.state.uploadMessage,uploadProgress:this.store.state.uploadProgress}},supporter:function(){return!!this.store.state.sponsor}},props:{notifications:Array},created:function(){const t=new URLSearchParams(window.location.search);this.compact=t.get("compact"),this.connect()}},vt=pt,ht=(a("8aef"),Object(A["a"])(vt,n,o,!1,null,"a63c8d36",null)),ft=ht.exports,gt=a("8c4f"),bt=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"container"},[t.configured?a("Site",t._b({},"Site",t.state,!1)):a("div",[a("div",{staticClass:"row py-5"},[a("div",{staticClass:"col12"},[a("p",{staticClass:"h1 pt-5 pb-2 border-bottom"},[t._v("Willkommen bei evcc")]),a("p",{staticClass:"lead pt-2"},[a("b",[t._v("evcc")]),t._v(" ist dient zur flexiblen Ladesteuerung von Elektrofahrzeugen. ")]),a("p",{staticClass:"pt-2"},[t._v(" Es sieht aus, als wäre Dein "),a("b",[t._v("evcc")]),t._v(" noch nicht konfiguriert. Um "),a("b",[t._v("evcc")]),t._v(" zu konfigurieren sind die folgenden Schritte notwendig: ")]),a("ol",{staticClass:"pt-2"},[a("li",[t._v(" Erzeugen einer Konfigurationsdatei mit Namen "),a("code",[t._v("evcc.yaml")]),t._v(". Die Standardkonfiguration "),a("code",[t._v("evcc.dist.yaml")]),t._v(" kann dafür als Vorlage dienen ("),a("a",{attrs:{href:"https://github.com/andig/evcc/blob/master/evcc.dist.yaml"}},[t._v("Download")]),t._v("). ")]),a("li",[t._v("Konfiguration der Wallbox als "),a("code",[t._v("chargers")]),t._v(".")]),a("li",[t._v(" Konfiguration des EVU Zählers und evtl. weiterer Zähler unter "),a("code",[t._v("meters")]),t._v(". ")]),a("li",[t._v(" Konfiguration des Netzanschlusses unter "),a("code",[t._v("site")]),t._v(". In einer Site wird der Netzanschluss mit dem konfigurierten EVU Zähler ("),a("code",[t._v("meter")]),t._v(") verbunden. ")]),a("li",[t._v(" Konfiguration eines Ladepunktes unter "),a("code",[t._v("loadpoints")]),t._v(". In einem Ladepunkt wird die konfigurierte Wallbox ("),a("code",[t._v("charger")]),t._v(") mit dem Ladepunkt verbunden. ")]),a("li",[t._v(" Start von "),a("b",[t._v("evcc")]),t._v(" mit der neu erstellten Konfiguration: "),a("code",[t._v("evcc -c evcc.yaml")])])]),a("p",[t._v("Minimale Beispielkonfiguration für "),a("b",[t._v("evcc")]),t._v(":")]),a("p",[a("code",[a("pre",{staticClass:"mx-3"},[t._v("                uri: localhost:7070 # Adresse für UI\n                interval: 10s # Regelintervall\n                meters:\n                - name: evu-zähler\n                type: ... # Detailkonfiguration des EVU Zählers\n                - name: ladezähler\n                type: ... # Detailkonfiguration des Ladezählers (optional)\n                chargers:\n                - name: wallbox\n                type: ... # Detailkonfiguration der Wallbox\n                site:\n                  title: Home\n                  meters:\n                  grid: evu-zähler # EVU Zähler\n                loadpoints:\n                - title: Ladepunkt # ui display name\n                  charger: wallbox # charger\n                  meters:\n                    charge: ladezähler # Ladezählers (optional)\n              ")])])]),a("p",[t._v(" Viel Spass mit "),a("b",[t._v("evcc")]),t._v("! Bei Problemen kannst Du uns auf "),a("a",{attrs:{href:"https://github.com/andig/evcc/issues"}},[t._v("GitHub")]),t._v(" erreichen. ")])])])])],1)},Ct=[],_t=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"flex-grow-1 d-flex flex-column"},[a("div",{staticClass:"row mt-4 pt-2"},[a("div",{staticClass:"d-none d-md-flex col-12 col-md-3 col-lg-4 align-items-end"},[a("p",{staticClass:"h1 text-truncate"},[t._v(t._s(t.siteTitle||"Home"))])]),a("div",{staticClass:"col-12 col-md-9 col-lg-6 flex-grow-1"},[a("SiteDetails",t._b({},"SiteDetails",t.details,!1))],1)]),a("hr",{staticClass:"w-100 my-4"}),a("div",{staticClass:"flex-grow-1 d-flex justify-content-around flex-column"},[t._l(t.loadpoints,(function(e,s){return[s>0?a("hr",{key:s+"_hr",staticClass:"w-100 my-4"}):t._e(),a("Loadpoint",t._b({key:s,attrs:{single:1===t.loadpoints.length,id:s,pvConfigured:t.pvConfigured}},"Loadpoint",e,!1))]}))],2)])},yt=[],wt=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"row row-cols-3 justify-content-between justify-content-md-start overflow-hidden"},[t.gridConfigured?a("div",{staticClass:"px-3"},[t.gridPower>0?a("div",{staticClass:"mb-2 value"},[t._v(" Bezug "),a("fa-icon",{staticClass:"text-primary",attrs:{icon:"arrow-down"}})],1):a("div",{staticClass:"mb-2 value"},[t._v(" Einspeisung "),a("fa-icon",{staticClass:"text-primary",attrs:{icon:"arrow-up"}})],1),a("h3",{staticClass:"value"},[t._v(" "+t._s(t.fmt(t.gridPower))+" "),a("small",{staticClass:"text-muted"},[t._v(t._s(t.fmtUnit(t.gridPower))+"W")])])]):t._e(),t.pvConfigured?a("div",{staticClass:"px-3"},[a("div",{staticClass:"mb-2 value"},[t._v(" Erzeugung "),a("fa-icon",{class:{"text-primary":t.pvPower>0,"text-muted":t.pvPower<=0},attrs:{icon:"sun"}})],1),a("h3",{staticClass:"value"},[t._v(" "+t._s(t.fmt(t.pvPower))+" "),a("small",{staticClass:"text-muted"},[t._v(t._s(t.fmtUnit(t.pvPower))+"W")])])]):t._e(),t.batteryConfigured?a("div",{staticClass:"px-3"},[a("div",{staticClass:"mb-2 value"},[a("div",{staticClass:"d-block d-sm-none"},[t._v(" Akku "),a("span",{staticClass:"text-muted"},[t._v(" / "+t._s(t.batterySoC)+" %")])]),a("div",{staticClass:"d-none d-sm-block"},[t._v(" Batterie "),a("span",{staticClass:"text-muted"},[t._v(" / "+t._s(t.batterySoC)+"% ")]),a("fa-icon",{staticClass:"text-primary",attrs:{icon:t.batteryIcon}})],1)]),a("h3",{staticClass:"value"},[t._v(" "+t._s(t.fmt(t.batteryPower))+" "),a("small",{staticClass:"text-muted"},[t._v(t._s(t.fmtUnit(t.batteryPower))+"W")])])]):t._e()])},St=[];const xt=20,kt=["battery-empty","battery-quarter","battery-half","battery-three-quarters","battery-full"];var Nt={name:"SiteDetails",props:{gridConfigured:Boolean,gridPower:Number,pvConfigured:Boolean,pvPower:Number,batteryConfigured:Boolean,batteryPower:Number,batterySoC:Number},data:function(){return{iconIdx:0}},mixins:[rt],computed:{numberOfPanels:function(){let t=0;return this.gridConfigured&&t++,this.pvConfigured&&t++,this.batteryConfigured&&t++,t},batteryIcon:function(){return Math.abs(this.batteryPower)<xt?this.batterySoC<30?kt[0]:this.batterySoC<50?kt[1]:this.batterySoC<70?kt[2]:this.batterySoC<90?kt[3]:kt[4]:kt[this.iconIdx]}},mounted:function(){window.setInterval(()=>{this.batteryPower>xt?--this.iconIdx<0&&(this.iconIdx=kt.length-1):this.batteryPower<xt&&++this.iconIdx>=kt.length&&(this.iconIdx=0)},1e3)}},Dt=Nt,Mt=Object(A["a"])(Dt,wt,St,!1,null,null,null),Pt=Mt.exports,Tt=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("p",{staticClass:"h3 mb-4 d-sm-block",class:{"d-none":t.single}},[t._v(t._s(t.title||"Ladepunkt"))]),"soft"==t.remoteDisabled?a("div",{staticClass:"alert alert-warning mt-4 mb-2",attrs:{role:"alert"}},[t._v(" "+t._s(t.remoteDisabledSource)+": Adaptives PV-Laden deaktiviert ")]):t._e(),"hard"==t.remoteDisabled?a("div",{staticClass:"alert alert-danger mt-4 mb-2",attrs:{role:"alert"}},[t._v(" "+t._s(t.remoteDisabledSource)+": Deaktiviert ")]):t._e(),a("div",{staticClass:"row"},[a("Mode",{staticClass:"col-12 col-md-6 col-lg-4 mb-4",attrs:{mode:t.mode,pvConfigured:t.pvConfigured},on:{updated:t.setTargetMode}}),a("Vehicle",t._b({staticClass:"col-12 col-md-6 col-lg-8 mb-4",on:{"target-soc-updated":t.setTargetSoC}},"Vehicle",t.vehicle,!1))],1),a("LoadpointDetails",t._b({},"LoadpointDetails",t.details,!1))],1)},Bt=[],Et=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("div",{staticClass:"mb-3"},[t._v("Modus")]),a("div",{staticClass:"btn-group w-100",attrs:{role:"group"}},[a("input",{staticClass:"btn-check",attrs:{id:"mode_off",type:"radio",name:"mode",value:"off"},domProps:{checked:"off"==t.mode},on:{click:function(e){return t.setTargetMode("off")}}}),a("label",{staticClass:"btn btn-outline-primary",attrs:{for:"mode_off"}},[t._v(" Stop ")]),a("input",{staticClass:"btn-check",attrs:{id:"mode_now",type:"radio",name:"mode",value:"now"},domProps:{checked:"now"==t.mode},on:{click:function(e){return t.setTargetMode("now")}}}),a("label",{staticClass:"btn btn-outline-primary",attrs:{for:"mode_now"}},[t._v(" Sofort ")]),a("input",{staticClass:"btn-check",attrs:{id:"mode_minpv",type:"radio",name:"mode",value:"minpv"},domProps:{checked:"minpv"==t.mode},on:{click:function(e){return t.setTargetMode("minpv")}}}),t.pvConfigured?a("label",{staticClass:"btn btn-outline-primary",attrs:{for:"mode_minpv"}},[a("span",{staticClass:"d-inline d-sm-none"},[t._v("Min")]),a("span",{staticClass:"d-none d-sm-inline"},[t._v("Min + PV")])]):t._e(),a("input",{staticClass:"btn-check",attrs:{id:"mode_pv",type:"radio",name:"mode",value:"pv"},domProps:{checked:"pv"==t.mode},on:{click:function(e){return t.setTargetMode("pv")}}}),t.pvConfigured?a("label",{staticClass:"btn btn-outline-primary",attrs:{for:"mode_pv"}},[a("span",{staticClass:"d-inline d-sm-none"},[t._v("PV")]),a("span",{staticClass:"d-none d-sm-inline"},[t._v("Nur PV")])]):t._e()])])},Vt=[],jt={name:"Mode",props:{mode:String,pvConfigured:Boolean},methods:{setTargetMode:function(t){this.$emit("updated",t)}}},At=jt,$t=(a("3c05"),Object(A["a"])(At,Et,Vt,!1,null,"685c2a71",null)),zt=$t.exports,Lt=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("div",{staticClass:"mb-3"},[t._v(" "+t._s(t.socTitle||"Fahrzeug")+" ")]),a("VehicleSoc",t._b({on:{"target-soc-updated":t.targetSocUpdated}},"VehicleSoc",t.vehicleSoc,!1)),a("VehicleSubline",t._b({staticClass:"my-1"},"VehicleSubline",t.vehicleSubline,!1))],1)},Ot=[],Ut={methods:{collectProps:function(t){let e={};for(var a in t.props)e[a]=this[a];return e}}},It=function(){var t,e,a=this,s=a.$createElement,i=a._self._c||s;return i("div",{staticClass:"vehicle-soc"},[i("div",{staticClass:"progress"},[i("div",{staticClass:"progress-bar",class:(t={"progress-bar-striped":a.charging,"progress-bar-animated":a.charging},t[a.progressColor]=!0,t),style:{width:a.socChargeDisplayWidth+"%"},attrs:{role:"progressbar"}},[a._v(" "+a._s(a.socChargeDisplayValue)+" ")]),a.remainingSoCWidth>0?i("div",{staticClass:"progress-bar",class:(e={},e[a.progressColor]=!0,e["bg-muted"]=!0,e),style:{width:a.remainingSoCWidth+"%",transition:"none"},attrs:{role:"progressbar"}}):a._e()]),a.connected&&a.hasVehicle&&a.visibleTargetSoC?i("div",{staticClass:"target",class:{"target--max":100===a.visibleTargetSoC}},[i("div",{staticClass:"target-label d-flex align-items-center justify-content-center",style:{left:a.visibleTargetSoC+"%"}},[a._v(" "+a._s(a.visibleTargetSoC)+"% ")]),i("input",{staticClass:"target-slider",attrs:{type:"range",min:"0",max:"100",step:"5"},domProps:{value:a.visibleTargetSoC},on:{input:a.movedTargetSoC,change:a.setTargetSoC}})]):a._e()])},Wt=[],Ht={name:"VehicleSoc",props:{connected:Boolean,hasVehicle:Boolean,socCharge:Number,enabled:Boolean,charging:Boolean,minSoC:Number,targetSoC:Number},data:function(){return{selectedTargetSoC:null}},computed:{socChargeDisplayWidth:function(){return this.hasVehicle&&this.socCharge>=0?this.socCharge:100},socChargeDisplayValue:function(){if(!this.hasVehicle||!this.socCharge||this.socCharge<0){let t="getrennt";return this.charging?t="lädt":this.enabled?t="bereit":this.connected&&(t="verbunden"),t}let t=this.socCharge;return t>=10&&(t+="%"),t},progressColor:function(){return this.connected?this.minSoCActive?"bg-danger":this.enabled?"bg-primary":"bg-secondary":"bg-light border"},minSoCActive:function(){return this.minSoC>0&&this.socCharge<this.minSoC},remainingSoCWidth:function(){return 100===this.socChargeDisplayWidth?null:this.minSoCActive?this.minSoC-this.socCharge:this.visibleTargetSoC>this.socCharge?this.visibleTargetSoC-this.socCharge:null},visibleTargetSoC:function(){return Number(this.selectedTargetSoC||this.targetSoC)}},methods:{movedTargetSoC:function(t){const e=40;return t.target.value<e?(t.target.value=e,this.selectedTargetSoC=t.target.value,t.preventDefault(),!1):(this.selectedTargetSoC=t.target.value,!0)},setTargetSoC:function(t){this.$emit("target-soc-updated",t.target.value)}}},Ft=Ht,Kt=(a("f673"),Object(A["a"])(Ft,It,Wt,!1,null,"61c611fd",null)),Rt=Kt.exports,Zt=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"d-flex justify-content-between align-items-center"},[a("small",{staticClass:"text-secondary"},[t.minSoCActive?a("span",[a("fa-icon",{staticClass:"text-muted me-1",attrs:{icon:"exclamation-circle"}}),t._v(" Mindestladung bis "+t._s(t.minSoC)+"% ")],1):t._e()]),t.targetChargeEnabled?a("small",{class:{invisible:!t.targetSoC,"text-primary":t.timerActive,"text-secondary":!t.timerActive}},[t._v(" "+t._s(t.targetTimeLabel())+" "),a("fa-icon",{staticClass:"ms-1",attrs:{icon:"clock"}})],1):t._e()])},Jt=[],Gt={name:"VehicleSubline",props:{socCharge:Number,minSoC:Number,timerActive:Boolean,timerSet:Boolean,targetTime:String,targetSoC:Number},computed:{minSoCActive:function(){return this.minSoC>0&&this.socCharge<this.minSoC},targetChargeEnabled:function(){return this.targetTime&&this.timerSet}},methods:{targetTimeLabel:function(){const t=new Date(this.targetTime);return`bis ${this.fmtAbsoluteDate(t)} Uhr`}},mixins:[rt]},qt=Gt,Qt=Object(A["a"])(qt,Zt,Jt,!1,null,null,null),Yt=Qt.exports,Xt={name:"Vehicle",components:{VehicleSoc:Rt,VehicleSubline:Yt},props:{connected:Boolean,hasVehicle:Boolean,socCharge:Number,enabled:Boolean,charging:Boolean,minSoC:Number,socTitle:String,timerActive:Boolean,timerSet:Boolean,targetTime:String,targetSoC:Number},computed:{vehicleSoc:function(){return this.collectProps(Rt)},vehicleSubline:function(){return this.collectProps(Yt)}},methods:{targetSocUpdated:function(t){this.$emit("target-soc-updated",t)}},mixins:[Ut]},te=Xt,ee=Object(A["a"])(te,Lt,Ot,!1,null,null,null),ae=ee.exports,se=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("div",{staticClass:"row"},[a("div",{staticClass:"col-6 col-sm-3 col-lg-2 mt-3 offset-lg-4"},[a("div",{staticClass:"mb-2 value"},[t._v(" Leistung "),"heating"==t.climater?a("fa-icon",{staticClass:"text-primary ms-1",attrs:{icon:"temperature-low"}}):t._e(),"cooling"==t.climater?a("fa-icon",{staticClass:"text-primary ms-1",attrs:{icon:"temperature-high"}}):t._e(),"on"==t.climater?a("fa-icon",{staticClass:"text-primary ms-1",attrs:{icon:"thermometer-half"}}):t._e()],1),a("h3",{staticClass:"value"},[t._v(" "+t._s(t.fmt(t.chargePower))+" "),a("small",{staticClass:"text-muted"},[t._v(t._s(t.fmtUnit(t.chargePower))+"W")])])]),a("div",{staticClass:"col-6 col-sm-3 col-lg-2 mt-3"},[a("div",{staticClass:"mb-2 value"},[t._v("Geladen")]),a("h3",{staticClass:"value"},[t._v(" "+t._s(t.fmt(t.chargedEnergy))+" "),a("small",{staticClass:"text-muted"},[t._v(t._s(t.fmtUnit(t.chargedEnergy))+"Wh")])])]),t.range>=0?a("div",{staticClass:"col-6 col-sm-3 col-lg-2 mt-3"},[a("div",{staticClass:"mb-2 value"},[t._v("Reichweite")]),a("h3",{staticClass:"value"},[t._v(" "+t._s(Math.round(t.range))+" "),a("small",{staticClass:"text-muted"},[t._v("km")])])]):a("div",{staticClass:"col-6 col-sm-3 col-lg-2 mt-3"},[a("div",{staticClass:"mb-2 value"},[t._v("Dauer")]),a("h3",{staticClass:"value"},[t._v(" "+t._s(t.fmtShortDuration(t.chargeDuration))+" "),a("small",{staticClass:"text-muted"},[t._v(t._s(t.fmtShortDurationUnit(t.chargeDuration)))])])]),t.hasVehicle?a("div",{staticClass:"col-6 col-sm-3 col-lg-2 mt-3"},[a("div",{staticClass:"mb-2 value"},[t._v("Restzeit")]),a("h3",{staticClass:"value"},[t._v(" "+t._s(t.fmtShortDuration(t.chargeEstimate))+" "),a("small",{staticClass:"text-muted"},[t._v(t._s(t.fmtShortDurationUnit(t.chargeEstimate)))])])]):t._e()])])},ie=[],re={name:"LoadpointDetails",props:{chargedEnergy:Number,chargeDuration:Number,chargeEstimate:Number,chargePower:Number,climater:String,hasVehicle:Boolean,range:Number},mixins:[rt]},ne=re,oe=Object(A["a"])(ne,se,ie,!1,null,null,null),le=oe.exports,ce={name:"Loadpoint",props:{id:Number,pvConfigured:Boolean,single:Boolean,title:String,mode:String,targetSoC:Number,remoteDisabled:Boolean,remoteDisabledSource:String,chargeDuration:Number,charging:Boolean,connected:Boolean,enabled:Boolean,socTitle:String,socCharge:Number,minSoC:Number,timerSet:Boolean,timerActive:Boolean,targetTime:String,chargePower:Number,chargedEnergy:Number,hasVehicle:Boolean,climater:String,range:Number,chargeEstimate:Number,phases:Number,minCurrent:Number,maxCurrent:Number,activePhases:Number,chargeCurrent:Number,socCapacity:Number,connectedDuration:Number,chargeCurrents:Array,chargeConfigured:Boolean,chargeRemainingEnergy:Number},components:{LoadpointDetails:le,Mode:zt,Vehicle:ae},mixins:[rt,Ut],data:function(){return{tickerHandle:null,chargeDurationDisplayed:null}},computed:{details:function(){return this.collectProps(le)},vehicle:function(){return this.collectProps(ae)}},watch:{chargeDuration:function(){window.clearInterval(this.tickerHandle),this.charging&&this.chargeDuration>=0&&(this.chargeDurationDisplayed=this.chargeDuration,this.tickerHandle=window.setInterval(function(){this.chargeDurationDisplayed+=1}.bind(this),1e3))}},methods:{api:function(t){return"loadpoints/"+this.id+"/"+t},setTargetMode:function(t){r.a.post(this.api("mode")+"/"+t).then(function(t){this.mode=t.data.mode}.bind(this)).catch(window.app.error)},setTargetSoC:function(t){r.a.post(this.api("targetsoc")+"/"+t).then(function(t){this.targetSoC=t.data.targetSoC}.bind(this)).catch(window.app.error)}},destroyed:function(){window.clearInterval(this.tickerHandle)}},de=ce,ue=Object(A["a"])(de,Tt,Bt,!1,null,null,null),me=ue.exports,pe={name:"Site",props:{siteTitle:String,loadpoints:Array,gridConfigured:Boolean,gridPower:Number,pvConfigured:Boolean,pvPower:Number,batteryConfigured:Boolean,batteryPower:Number,batterySoC:Number,gridCurrents:Array,prioritySoC:Number},components:{SiteDetails:Pt,Loadpoint:me},mixins:[rt,Ut],computed:{details:function(){return this.collectProps(Pt)}}},ve=pe,he=Object(A["a"])(ve,_t,yt,!1,null,null,null),fe=he.exports,ge={name:"Main",components:{Site:fe},data:function(){return this.$root.$data.store},computed:{configured:function(){const t=window.evcc.configured;return t==window.evcc.configured||!isNaN(parseInt(t))&&parseInt(t)>0}}},be=ge,Ce=Object(A["a"])(be,bt,Ct,!1,null,null,null),_e=Ce.exports,ye=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"container"},[a("div",{staticClass:"row mt-4 border-bottom"},[a("div",{staticClass:"col-12"},[a("p",{staticClass:"h1"},[t._v(t._s(t.title||"Home"))])])]),a("div",{staticClass:"row h5"},[a("div",{staticClass:"col-md-4"}),a("div",{staticClass:"col-6 col-md-2 py-3"},[t._v(" Netzzähler: "),t.gridConfigured?a("span",{staticClass:"text-primary"},[t._v("✓")]):a("span",{staticClass:"text-primary"},[t._v("—")])]),a("div",{staticClass:"col-6 col-md-2 py-3"},[t._v(" PV Zähler: "),t.pvConfigured?a("span",{staticClass:"text-primary"},[t._v("✓")]):a("span",{staticClass:"text-primary"},[t._v("—")])]),a("div",{staticClass:"col-6 col-md-2 py-3"},[t._v(" Batteriezähler: "),t.batteryConfigured?a("span",{staticClass:"text-primary"},[t._v("✓")]):a("span",{staticClass:"text-primary"},[t._v("—")])])]),t._l(t.loadpoints,(function(e,s){return a("div",{key:s,attrs:{loadpoint:e,id:"loadpoint-"+s}},[a("div",{staticClass:"row mt-4 border-bottom"},[a("div",{staticClass:"col-12"},[a("p",{staticClass:"h1"},[t._v(t._s(e.title||"Ladepunkt"))])])]),a("div",{staticClass:"row h5"},[a("div",{staticClass:"col-md-4"}),a("div",{staticClass:"col-6 col-md-2 py-3"},[t._v(" Ladezähler: "),e.chargeConfigured?a("span",{staticClass:"text-primary"},[t._v("✓")]):a("span",{staticClass:"text-primary"},[t._v("—")])]),a("div",{staticClass:"col-6 col-md-2 py-3"},[t._v(" Phasen: "),a("span",{staticClass:"text-primary"},[t._v(t._s(e.phases)+"p")])]),a("div",{staticClass:"col-6 col-md-2 py-3"},[t._v(" Min. Strom: "),a("span",{staticClass:"text-primary"},[t._v(t._s(e.minCurrent)+"A")])]),a("div",{staticClass:"col-6 col-md-2 py-3"},[t._v(" Max. Strom: "),a("span",{staticClass:"text-primary"},[t._v(t._s(e.maxCurrent)+"A")])])]),a("div",{staticClass:"row h5"},[a("div",{staticClass:"col-md-4"}),a("div",{staticClass:"col-md-8 h2"},[t._m(0,!0),a("div",{staticClass:"row h5"},[a("div",{staticClass:"col-6 py-3"},[t._v(" Modell: "),a("span",{staticClass:"text-primary"},[t._v(t._s(e.socTitle||"—"))])]),a("div",{staticClass:"col-6 py-3"},[t._v(" Kapazität: "),a("span",{staticClass:"text-primary"},[t._v(t._s(e.socCapacity)+"kWh")])])])])])])}))],2)},we=[function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"row py-3 h2 border-bottom"},[a("div",{staticClass:"col-12"},[t._v("Fahrzeug")])])}],Se={name:"Config",data:function(){return this.$root.$data.store.state}},xe=Se,ke=Object(A["a"])(xe,ye,we,!1,null,null,null),Ne=ke.exports;s["a"].use(gt["a"]);const De=[{path:"/",component:_e},{path:"/config",component:Ne}];var Me=new gt["a"]({routes:De,linkExactActiveClass:"active"});const Pe=window.location;r.a.defaults.baseURL=Pe.protocol+"//"+Pe.hostname+(Pe.port?":"+Pe.port:"")+Pe.pathname+"api",r.a.defaults.headers.post["Content-Type"]="application/json",window.app=new s["a"]({el:"#app",router:Me,data:{store:mt,notifications:[]},render:function(t){return t(ft,{props:{notifications:this.notifications}})},methods:{raise:function(t){console[t.type](t);const e=this.notifications.filter(e=>e.message!==t.message);this.notifications=[t,...e]},clear:function(){this.notifications=[]},error:function(t){t.type="error",this.raise(t)},warn:function(t){t.type="warn",this.raise(t)}}}),window.setInterval((function(){r.a.get("health").catch((function(){window.app.error({message:"Server unavailable"})}))}),5e3)},b19f:function(t,e,a){},c0cc:function(t,e,a){},f673:function(t,e,a){"use strict";a("2110")}});
//# sourceMappingURL=index.6f1d74b5.js.map