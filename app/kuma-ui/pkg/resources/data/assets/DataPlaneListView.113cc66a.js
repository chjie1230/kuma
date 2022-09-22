import{L as se,O as C,o as i,f as d,h as a,t as y,p as ne,M as b,_ as q,d as A,w as I,e as N,F as k,j as O,b as L,n as K,cx as oe,a2 as ee,cy as le,cz as we,r as E,z as ke,D as ie,E as re,q as Pe,cA as Te,cB as Ce,a3 as Ee,k as Oe,u as Se,x as Ie,a0 as te,cC as Ke,cD as Ae,cE as Ve,a5 as Ue,a6 as Le,cF as Ne,a4 as Me,K as S,i as Re,cG as $e}from"./index.2f6d90b0.js";import{g as Be}from"./tableDataUtils.f37d211d.js";import{D as xe}from"./DataOverview.cd1778b1.js";import{Y as ze}from"./YamlView.25baae1a.js";import{_ as He}from"./ErrorBlock.65239dfd.js";import"./CodeBlock.45d0fc63.js";import"./index.58caa11d.js";const ce=[{key:"status",label:"Status"},{key:"name",label:"Name"},{key:"mesh",label:"Mesh"},{key:"type",label:"Type"},{key:"service",label:"Service"},{key:"protocol",label:"Protocol"},{key:"zone",label:"Zone"},{key:"lastConnected",label:"Last Connected"},{key:"lastUpdated",label:"Last Updated"},{key:"totalUpdates",label:"Total Updates"},{key:"dpVersion",label:"Kuma DP version"},{key:"envoyVersion",label:"Envoy version"},{key:"details",label:"Details",hideLabel:!0}],je=["name","details"],qe=ce.filter(e=>!je.includes(e.key)).map(e=>({tableHeaderKey:e.key,label:e.label,isChecked:!1})),de=["status","name","mesh","type","service","protocol","zone","lastUpdated","dpVersion","details"];function Fe(e,t=de){return ce.filter(o=>t.includes(o.key)?e?!0:o.key!=="zone":!1)}const ae={get(e){try{const t=window.localStorage.getItem(e);return t!==null?JSON.parse(t):null}catch{return null}},set(e,t){try{const o=JSON.stringify(t);window.localStorage.setItem(e,o)}catch{}},remove(e){try{window.localStorage.removeItem(e)}catch{}}};function j(e,t){const o=window.history.state;if(o===null)return;const u=o.current.indexOf("?"),n=u>-1?o.current.substring(u):"",c=new URLSearchParams(n);t!=null?c.set(e,String(t)):c.has(e)&&c.delete(e);const f=c.toString(),D=f===""?"":"?"+f;let g="";if(u>-1?g=o.current.substring(0,u)+D:g=o.current+D,o.current!==g){const l=Object.assign({},o);l.current=g,window.history.replaceState(l,"","#"+g)}}const Ge={class:"entity-tag"},Qe={class:"entity-tag__value"},Ye=se({__name:"EntityTag",props:{tag:{type:Object,required:!0}},setup(e){const t=e,o=C(()=>t.tag.label.toLowerCase().includes("kuma.io/"));return(u,n)=>(i(),d("span",Ge,[a("span",{class:ne(["entity-tag__label",{"entity-tag__label--is-kuma-io-label":b(o)}])},y(e.tag.label),3),a("span",Qe,y(e.tag.value),1)]))}});const Ze=q(Ye,[["__scopeId","data-v-d4e000a4"]]),P=e=>(ie("data-v-532cb754"),e=e(),re(),e),Je={class:"entity-summary entity-section-list"},Xe={class:"entity-title"},We=P(()=>a("span",{class:"kutil-sr-only"},"Data plane:",-1)),et={class:"definition"},tt=P(()=>a("span",null,"Mesh:",-1)),at={key:0},st=P(()=>a("h4",null,"Tags",-1)),nt={class:"entity-tag-list"},ot={key:1},lt=P(()=>a("h4",null,"Dependencies",-1)),it={class:"mt-2 heading-with-icon"},rt=N(" Warnings "),ct=P(()=>a("h4",null,"Insights",-1)),dt={class:"entity-section-list"},pt=["data-testid"],ut=P(()=>a("span",null,"Connect time:",-1)),mt=["data-testid"],ht=P(()=>a("span",null,"Disconnect time:",-1)),yt={class:"definition"},vt=P(()=>a("span",null,"Control plane instance ID:",-1)),_t={key:0},gt=P(()=>a("summary",null," Responses (acknowledged / sent) ",-1)),bt=["data-testid"],ft=se({__name:"DataPlaneEntitySummary",props:{dataPlaneOverview:{type:Object,required:!0}},setup(e){const t=e,o={Online:"success",Offline:"danger","Partially degraded":"warning"},u=C(()=>{const{name:l,mesh:p,dataplane:v}=t.dataPlaneOverview;return{type:"Dataplane",name:l,mesh:p,networking:v.networking}}),n=C(()=>oe(t.dataPlaneOverview.dataplane)),c=C(()=>{const l=Array.from(t.dataPlaneOverview.dataplaneInsight.subscriptions);return l.reverse(),l.map(p=>{const v=p.connectTime!==void 0?ee(p.connectTime):"\u2014",r=p.disconnectTime!==void 0?ee(p.disconnectTime):"\u2014",s=Object.entries(p.status).filter(([m])=>!["total","lastUpdateTime"].includes(m)).map(([m,_])=>{var M,V,U,R,w;const T=`${(M=_.responsesAcknowledged)!=null?M:0} / ${(V=_.responsesSent)!=null?V:0}`;return{type:m.toUpperCase(),ratio:T,responsesSent:(U=_.responsesSent)!=null?U:0,responsesAcknowledged:(R=_.responsesAcknowledged)!=null?R:0,responsesRejected:(w=_.responsesRejected)!=null?w:0}});return{subscription:p,formattedConnectDate:v,formattedDisconnectDate:r,statuses:s}})}),f=C(()=>{const{status:l}=le(t.dataPlaneOverview.dataplane,t.dataPlaneOverview.dataplaneInsight);return{title:l,appearance:o[l]}}),D=C(()=>{const l=we(t.dataPlaneOverview.dataplaneInsight);return l!==null?Object.entries(l).map(([p,v])=>({name:p,version:v})):[]}),g=C(()=>{const{subscriptions:l}=t.dataPlaneOverview.dataplaneInsight;if(l.length===0)return[];const p=[],v=l[l.length-1],r=v.version.envoy,s=v.version.kumaDp,m=r.kumaDpCompatible!==void 0?r.kumaDpCompatible:!0,_=s.kumaCpCompatible!==void 0?s.kumaCpCompatible:!0;if(!m){const T=`Envoy ${r.version} is not supported by Kuma DP ${s.version}.`;p.push(T)}if(!_){const T=`Kuma DP ${s.version} is not supported by this Kuma control plane.`;p.push(T)}return p});return(l,p)=>{const v=E("router-link");return i(),d("div",Je,[a("section",null,[a("h3",Xe,[We,A(v,{to:{name:"data-plane-detail-view",params:{mesh:e.dataPlaneOverview.mesh,dataPlane:e.dataPlaneOverview.name}}},{default:I(()=>[N(y(e.dataPlaneOverview.name),1)]),_:1},8,["to"]),a("div",{class:ne(`status status--${b(f).appearance}`),"data-testid":"data-plane-status-badge"},y(b(f).title.toLowerCase()),3)]),a("div",et,[tt,a("span",null,y(e.dataPlaneOverview.mesh),1)])]),b(n).length>0?(i(),d("section",at,[st,a("div",nt,[(i(!0),d(k,null,O(b(n),(r,s)=>(i(),L(Ze,{key:s,tag:r},null,8,["tag"]))),128))])])):K("",!0),b(D).length>0?(i(),d("section",ot,[lt,(i(!0),d(k,null,O(b(D),(r,s)=>(i(),d("div",{key:s,class:"definition"},[a("span",null,y(r.name)+":",1),a("span",null,y(r.version),1)]))),128)),b(g).length>0?(i(),d(k,{key:0},[a("h5",it,[rt,A(b(ke),{class:"ml-1",icon:"warning",color:"var(--black-75)","secondary-color":"var(--yellow-300)",size:"20"})]),(i(!0),d(k,null,O(b(g),(r,s)=>(i(),d("p",{key:s},y(r),1))),128))],64)):K("",!0)])):K("",!0),b(c).length>0?(i(),d(k,{key:2},[a("section",null,[ct,a("div",dt,[(i(!0),d(k,null,O(b(c),(r,s)=>(i(),d("div",{key:s},[a("div",{class:"definition","data-testid":`data-plane-connect-time-${s}`},[ut,a("span",null,y(r.formattedConnectDate),1)],8,pt),a("div",{class:"definition","data-testid":`data-plane-disconnect-time-${s}`},[ht,a("span",null,y(r.formattedDisconnectDate),1)],8,mt),a("div",yt,[vt,a("span",null,y(r.subscription.controlPlaneInstanceId),1)]),r.statuses.length>0?(i(),d("details",_t,[gt,(i(!0),d(k,null,O(r.statuses,(m,_)=>(i(),d("div",{key:`${s}-${_}`,class:"definition","data-testid":`data-plane-subscription-status-${s}-${_}`},[a("span",null,y(m.type)+":",1),a("span",null,y(m.ratio),1)],8,bt))),128))])):K("",!0)]))),128))])]),a("section",null,[A(ze,{content:b(u)},null,8,["content"])])],64)):K("",!0)])}}});const Dt=q(ft,[["__scopeId","data-v-532cb754"]]);const wt={name:"DataPlaneListView",dataPlaneTypes:["All","Standard","Gateway (builtin)","Gateway (provided)"],emptyStateMsg:"There are no data plane proxies present.",nsBackButtonRoute:{name:"data-plane-list-view"},dataplaneApiParams:{},components:{DataOverview:xe,DataPlaneEntitySummary:Dt,KButton:Pe,KDropdownItem:Te,KDropdownMenu:Ce,EmptyBlock:He},props:{name:{type:String,required:!1,default:null},offset:{type:Number,required:!1,default:0}},data(){return{visibleTableHeaderKeys:de,productName:Ee,isLoading:!0,isEmpty:!1,hasError:!1,tableDataIsEmpty:!1,tableData:{headers:[],data:[]},pageSize:50,next:null,shownTLSTab:!1,rawData:null,filteredDataPlaneType:"All",pageOffset:this.offset,dataPlaneOverview:null}},computed:{...Oe({environment:"config/getEnvironment",queryNamespace:"getItemQueryNamespace",multicluster:"config/getMulticlusterStatus"}),dataplaneWizardRoute(){return this.environment==="universal"?{name:"universal-dataplane"}:{name:"kubernetes-dataplane"}},filteredTableData(){const e=this.tableData.data.filter(o=>this.filteredDataPlaneType==="All"?!0:o.type.toLowerCase()===this.filteredDataPlaneType.toLowerCase()),t=Fe(this.multicluster,this.visibleTableHeaderKeys);return{data:e,headers:t}},columnsDropdownItems(){return qe.filter(e=>this.multicluster?!0:e.tableHeaderKey!=="zone").map(e=>{const t=this.visibleTableHeaderKeys.includes(e.tableHeaderKey);return{...e,isChecked:t}})}},watch:{"$route.params.mesh":function(){this.$route.name==="data-plane-list-view"&&(this.isLoading=!0,this.isEmpty=!1,this.hasError=!1,this.tableDataIsEmpty=!1,this.loadData(0))}},created(){const e=ae.get("dpVisibleTableHeaderKeys");Array.isArray(e)&&(this.visibleTableHeaderKeys=e)},beforeMount(){this.loadData(this.offset)},methods:{stopPropagatingClickEvent(e){e.stopPropagation()},updateVisibleTableHeaders(e,t){const o=e.target,u=this.visibleTableHeaderKeys.findIndex(n=>n===t);o.checked&&u===-1?this.visibleTableHeaderKeys.push(t):!o.checked&&u>-1&&this.visibleTableHeaderKeys.splice(u,1),ae.set("dpVisibleTableHeaderKeys",Array.from(new Set(this.visibleTableHeaderKeys)))},onCreateClick(){Se.logger.info(Ie.CREATE_DATA_PLANE_PROXY_CLICKED)},getEmptyState(){return{title:"No Data",message:this.$options.emptyStateMsg}},async parseData(e){var G,Q,Y,Z,J,X;const{dataplane:t={},dataplaneInsight:o={}}=e,{name:u="",mesh:n=""}=e,{subscriptions:c=[]}=o,f={name:"data-plane-detail-view",params:{mesh:n,dataPlane:u}},D={name:"mesh-child",params:{mesh:n}},g=["kuma.io/protocol","kuma.io/service","kuma.io/zone"],l=oe(t).filter(h=>g.includes(h.label)),p=(Q=(G=l.find(h=>h.label==="kuma.io/service"))==null?void 0:G.value)!=null?Q:"\u2014",v=(Z=(Y=l.find(h=>h.label==="kuma.io/protocol"))==null?void 0:Y.value)!=null?Z:"\u2014",r=(X=(J=l.find(h=>h.label==="kuma.io/zone"))==null?void 0:J.value)!=null?X:"\u2014",{status:s}=le(t,o),{totalUpdates:m,totalRejectedUpdates:_,dpVersion:T,envoyVersion:M,selectedTime:V,selectedUpdateTime:U,version:R}=c.reduce((h,pe)=>{const{status:ue={},connectTime:me,version:W={}}=pe,{total:he={},lastUpdateTime:ye}=ue,{responsesSent:ve="0",responsesRejected:_e="0"}=he,{kumaDp:ge={},envoy:be={}}=W,{version:fe}=ge,{version:De}=be;let{selectedTime:$,selectedUpdateTime:B}=h;const z=Date.parse(me),H=Date.parse(ye);return z&&(!$||z>$)&&($=z),H&&(!B||H>B)&&(B=H),{totalUpdates:h.totalUpdates+parseInt(ve,10),totalRejectedUpdates:h.totalRejectedUpdates+parseInt(_e,10),dpVersion:fe||h.dpVersion,envoyVersion:De||h.envoyVersion,selectedTime:$,selectedUpdateTime:B,version:W||h.version}},{totalUpdates:0,totalRejectedUpdates:0,dpVersion:"\u2014",envoyVersion:"\u2014",selectedTime:NaN,selectedUpdateTime:NaN,version:{}}),w={name:u,nameRoute:f,mesh:n,meshRoute:D,zone:r,service:p,protocol:v,status:s,totalUpdates:m,totalRejectedUpdates:_,dpVersion:T,envoyVersion:M,warnings:[],unsupportedEnvoyVersion:!1,unsupportedKumaDPVersion:!1,kumaDpAndKumaCpMismatch:!1,lastUpdated:U?te(new Date(U).toUTCString()):"\u2014",lastConnected:V?te(new Date(V).toUTCString()):"\u2014",type:Ke(t)},{kind:x}=Ae(R);switch(x!==Ve&&w.warnings.push(x),x){case Le:w.unsupportedEnvoyVersion=!0;break;case Ue:w.unsupportedKumaDPVersion=!0;break}if(this.multicluster){const{compatible:h}=await Ne(l,T);h||(w.warnings.push(Me),w.kumaDpAndKumaCpMismatch=!0)}return w},async loadData(e){var u;this.isLoading=!0,this.pageOffset=e,j("offset",e>0?e:null);const t=this.$route.params.mesh||null,o=this.$route.query.ns||null;try{const{data:n,next:c}=await Be({getSingleEntity:S.getDataplaneOverviewFromMesh.bind(S),getAllEntities:S.getAllDataplaneOverviews.bind(S),getAllEntitiesFromMesh:S.getAllDataplaneOverviewsFromMesh.bind(S),size:this.pageSize,offset:e,mesh:t,query:o,params:{...this.$options.dataplaneApiParams}});if(n.length>0){this.next=c,this.rawData=n,this.selectDataPlaneOverview((u=this.name)!=null?u:n[0].name);const f=await Promise.all(n.map(D=>this.parseData(D)));this.tableData.data=f,this.tableDataIsEmpty=!1,this.isEmpty=!1}else this.selectDataPlaneOverview(null),this.tableData.data=[],this.tableDataIsEmpty=!0,this.isEmpty=!0}catch(n){this.hasError=!0,this.isEmpty=!0,console.error(n)}finally{this.isLoading=!1}},async selectDataPlaneOverview(e){var t;e?(this.dataPlaneOverview=(t=this.rawData.find(o=>o.name===e))!=null?t:this.rawData[0],j("name",this.dataPlaneOverview.name)):(this.dataPlaneOverview=null,j("name",null))}}},F=e=>(ie("data-v-853caa4a"),e=e(),re(),e),kt={class:"data-planes-container"},Pt={class:"data-planes-content component-frame"},Tt=F(()=>a("label",{for:"data-planes-type-filter",class:"mr-2"}," Type: ",-1)),Ct=["value"],Et=["for"],Ot=["id","checked","onChange"],St=F(()=>a("span",{class:"custom-control-icon"}," + ",-1)),It=N(" Create data plane proxy "),Kt=F(()=>a("span",{class:"custom-control-icon"}," \u2190 ",-1)),At=N(" View All "),Vt={class:"data-planes-sidebar component-frame"};function Ut(e,t,o,u,n,c){var r;const f=E("KDropdownItem"),D=E("KDropdownMenu"),g=E("KButton"),l=E("DataOverview"),p=E("DataPlaneEntitySummary"),v=E("EmptyBlock");return i(),d("div",kt,[a("div",Pt,[A(l,{"selected-entity-name":(r=n.dataPlaneOverview)==null?void 0:r.name,"page-size":n.pageSize,"has-error":n.hasError,"is-loading":n.isLoading,"empty-state":c.getEmptyState(),"table-data":c.filteredTableData,"table-data-is-empty":n.tableDataIsEmpty,"show-details":"",next:n.next,"page-offset":n.pageOffset,onTableAction:t[2]||(t[2]=s=>c.selectDataPlaneOverview(s.name)),onLoadData:t[3]||(t[3]=s=>c.loadData(s))},{additionalControls:I(()=>[a("div",null,[Tt,Re(a("select",{id:"data-planes-type-filter","onUpdate:modelValue":t[0]||(t[0]=s=>n.filteredDataPlaneType=s),"data-testid":"data-planes-type-filter"},[(i(!0),d(k,null,O(e.$options.dataPlaneTypes,(s,m)=>(i(),d("option",{key:m,value:s},y(s),9,Ct))),128))],512),[[$e,n.filteredDataPlaneType]])]),A(D,{label:"Columns",icon:"cogwheel","button-appearance":"outline"},{items:I(()=>[a("div",{onClick:t[1]||(t[1]=(...s)=>c.stopPropagatingClickEvent&&c.stopPropagatingClickEvent(...s))},[(i(!0),d(k,null,O(c.columnsDropdownItems,(s,m)=>(i(),L(f,{key:m,class:"table-header-selector-item",item:s},{default:I(()=>[a("label",{for:`data-plane-table-header-checkbox-${m}`,class:"k-checkbox table-header-selector-item-checkbox"},[a("input",{id:`data-plane-table-header-checkbox-${m}`,checked:s.isChecked,type:"checkbox",class:"k-input",onChange:_=>c.updateVisibleTableHeaders(_,s.tableHeaderKey)},null,40,Ot),N(" "+y(s.label),1)],8,Et)]),_:2},1032,["item"]))),128))])]),_:1}),A(g,{class:"add-dp-button",appearance:"primary",to:c.dataplaneWizardRoute,"data-testid":"data-plane-create-data-plane-button",onClick:c.onCreateClick},{default:I(()=>[St,It]),_:1},8,["to","onClick"]),e.$route.query.ns?(i(),L(g,{key:0,appearance:"primary",to:e.$options.nsBackButtonRoute,"data-testid":"data-plane-ns-back-button"},{default:I(()=>[Kt,At]),_:1},8,["to"])):K("",!0)]),_:1},8,["selected-entity-name","page-size","has-error","is-loading","empty-state","table-data","table-data-is-empty","next","page-offset"])]),a("div",Vt,[n.dataPlaneOverview!==null?(i(),L(p,{key:0,"data-plane-overview":n.dataPlaneOverview},null,8,["data-plane-overview"])):(i(),L(v,{key:1}))])])}const zt=q(wt,[["render",Ut],["__scopeId","data-v-853caa4a"]]);export{zt as default};