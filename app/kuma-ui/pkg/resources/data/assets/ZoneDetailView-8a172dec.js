import{d as C,c as _,Z as I,r as V,o as r,a as c,w as e,q as m,f as B,g as n,h as t,b as a,G as z,t as u,M as D,e as T,F as $,s as E}from"./index-bb115be3.js";import{g as A,r as Z,s as L,D as g,S as R,f as q,A as M,o as G,p as J,E as P,q as j,_ as F}from"./RouteView.vue_vue_type_script_setup_true_lang-e50bfc04.js";import{a as W,A as H,S as K,b as U}from"./SubscriptionHeader-9dfee5b9.js";import{_ as Q}from"./CodeBlock.vue_vue_type_style_index_0_lang-bc7de660.js";import{T as X}from"./TabsWidget-4ef098a4.js";import{T as Y}from"./TextWithCopyButton-8bd76c4b.js";import{_ as ee}from"./WarningsWidget.vue_vue_type_script_setup_true_lang-0f4368f4.js";import{f as te}from"./dataplane-30467516.js";import{_ as ne}from"./RouteTitle.vue_vue_type_script_setup_true_lang-c0548d11.js";import"./CopyButton-9e08aa37.js";const se={class:"stack"},oe={class:"variable-columns"},ie=C({__name:"ZoneDetails",props:{zoneOverview:{type:Object,required:!0}},setup(w){const o=w,{t:l}=A(),k=Z(),v=[{hash:"#overview",title:l("zone-cps.routes.item.tabs.overview")},{hash:"#insights",title:l("zone-cps.routes.item.tabs.insights")}],f=_(()=>{var i;for(const s of((i=o.zoneOverview.zoneInsight)==null?void 0:i.subscriptions)??[])if(s.config)return JSON.parse(s.config).environment;return"kubernetes"}),h=_(()=>L(o.zoneOverview)),b=_(()=>I(o.zoneOverview)),O=_(()=>{var p;const i=[],s=((p=o.zoneOverview.zoneInsight)==null?void 0:p.subscriptions)??[];if(s.length>0){const d=s[s.length-1],y=d.version.kumaCp.version||"-",{kumaCpGlobalCompatible:x=!0}=d.version.kumaCp;x||i.push({kind:te,payload:{zoneCpVersion:y,globalCpVersion:k("KUMA_VERSION")}})}return i}),S=_(()=>{var s;const i=((s=o.zoneOverview.zoneInsight)==null?void 0:s.subscriptions)??[];if(i.length>0){const p=i[i.length-1];if(p.config)return JSON.stringify(JSON.parse(p.config),null,2)}return null}),N=_(()=>{var s;const i=((s=o.zoneOverview.zoneInsight)==null?void 0:s.subscriptions)??[];return Array.from(i).reverse()});return(i,s)=>{const p=V("RouterLink");return r(),c(X,{tabs:v},{overview:e(()=>[m("div",se,[O.value.length>0?(r(),c(ee,{key:0,warnings:O.value},null,8,["warnings"])):B("",!0),n(),t(a(z),null,{body:e(()=>[m("div",oe,[t(g,null,{title:e(()=>[n(u(a(l)("http.api.property.status")),1)]),body:e(()=>[t(R,{status:h.value},null,8,["status"])]),_:1}),n(),t(g,null,{title:e(()=>[n(u(a(l)("http.api.property.name")),1)]),body:e(()=>[t(Y,{text:o.zoneOverview.name},{default:e(()=>[t(p,{to:{name:"zone-cp-detail-view",params:{zone:o.zoneOverview.name}}},{default:e(()=>[n(u(o.zoneOverview.name),1)]),_:1},8,["to"])]),_:1},8,["text"])]),_:1}),n(),t(g,null,{title:e(()=>[n(u(a(l)("http.api.property.type")),1)]),body:e(()=>[n(u(f.value),1)]),_:1}),n(),t(g,null,{title:e(()=>[n(u(a(l)("http.api.property.authenticationType")),1)]),body:e(()=>[n(u(b.value),1)]),_:1})])]),_:1}),n(),m("div",null,[m("h2",null,u(a(l)("zone-cps.detail.configuration_title")),1),n(),t(a(z),{class:"mt-4"},{body:e(()=>[S.value!==null?(r(),c(Q,{key:0,id:"code-block-zone-config",language:"json",code:S.value,"is-searchable":"","query-key":"zone-config"},null,8,["code"])):(r(),c(a(D),{key:1,class:"mt-4","data-testid":"warning-no-subscriptions",appearance:"warning"},{alertMessage:e(()=>[n(u(a(l)("zone-cps.detail.no_subscriptions")),1)]),_:1}))]),_:1})])])]),insights:e(()=>[t(a(z),null,{body:e(()=>[t(W,{"initially-open":0},{default:e(()=>[(r(!0),T($,null,E(N.value,(d,y)=>(r(),c(H,{key:y},{"accordion-header":e(()=>[t(K,{subscription:d},null,8,["subscription"])]),"accordion-content":e(()=>[t(U,{subscription:d},null,8,["subscription"])]),_:2},1024))),128))]),_:1})]),_:1})]),_:1})}}});const ae=q(ie,[["__scopeId","data-v-61561c57"]]),be=C({__name:"ZoneDetailView",setup(w){const{t:o}=A();return(l,k)=>(r(),c(F,{name:"zone-cp-detail-view","data-testid":"zone-cp-detail-view"},{default:e(({route:v})=>[t(M,{breadcrumbs:[{to:{name:"zone-cp-list-view"},text:a(o)("zone-cps.routes.item.breadcrumbs")}]},{title:e(()=>[m("h1",null,[t(ne,{title:a(o)("zone-cps.routes.item.title",{name:v.params.zone}),render:!0},null,8,["title"])])]),default:e(()=>[n(),t(G,{src:`/zone-cps/${v.params.zone}`},{default:e(({data:f,isLoading:h,error:b})=>[h?(r(),c(J,{key:0})):b!==void 0?(r(),c(P,{key:1,error:b},null,8,["error"])):f===void 0?(r(),c(j,{key:2})):(r(),c(ae,{key:3,"zone-overview":f,"data-testid":"detail-view-details"},null,8,["zone-overview"]))]),_:2},1032,["src"])]),_:2},1032,["breadcrumbs"])]),_:1}))}});export{be as default};
