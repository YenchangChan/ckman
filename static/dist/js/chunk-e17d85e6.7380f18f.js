(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-e17d85e6"],{"0148":function(t,e,n){"use strict";function a(t,e,n,a){var r=t.length,l=n+(a?1:-1);while(a?l--:++l<r)if(e(t[l],l,t))return l;return-1}var r=a;function l(t){return t!==t}var i=l;function c(t,e,n){var a=n-1,r=t.length;while(++a<r)if(t[a]===e)return a;return-1}var o=c;function u(t,e,n){return e===e?o(t,e,n):r(t,i,n)}e["a"]=u},"0a8e":function(t,e,n){"use strict";var a=n("a3cf"),r=n.n(a);r.a},"28e8":function(t,e,n){},3249:function(t,e,n){"use strict";var a=n("28e8"),r=n.n(a);r.a},"8cbb":function(t,e,n){"use strict";var a=n("9ac7"),r=n("0148");function l(t,e){var n=null==t?0:t.length;return!!n&&Object(r["a"])(t,e,0)>-1}var i=l;function c(t,e,n){var a=-1,r=null==t?0:t.length;while(++a<r)if(n(e,t[a]))return!0;return!1}var o=c,u=n("6568"),s=n("a55c");function f(){}var h=f,b=n("1989"),p=1/0,v=s["a"]&&1/Object(b["a"])(new s["a"]([,-0]))[1]==p?function(t){return new s["a"](t)}:h,d=v,m=200;function g(t,e,n){var r=-1,l=i,c=t.length,s=!0,f=[],h=f;if(n)s=!1,l=o;else if(c>=m){var p=e?null:d(t);if(p)return Object(b["a"])(p);s=!1,l=u["a"],h=new a["a"]}else h=e?[]:f;t:while(++r<c){var v=t[r],g=e?e(v):v;if(v=n||0!==v?v:0,s&&g===g){var y=h.length;while(y--)if(h[y]===g)continue t;e&&h.push(g),f.push(v)}else l(h,g,n)||(h!==f&&h.push(g),f.push(v))}return f}e["a"]=g},"8d70":function(t,e,n){"use strict";n.r(e);var a=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("main",{staticClass:"settings"},[n("breadcrumb",{attrs:{data:["Clusters",t.$route.params.id,"tables"]}}),n("section",[n("table-metric"),n("replication-table"),n("zk-table")],1)],1)},r=[],l=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"zkTable"},[n("div",{staticClass:"title flex flex-between flex-vcenter ptb-10"},[n("span",{staticClass:"fs-20 font-bold"},[t._v(t._s(t.$t("tables.Zookeeper Status")))])]),n("el-table",{staticClass:"tb-edit",staticStyle:{width:"100%"},attrs:{data:t.tableData,border:""}},t._l(t.cols,(function(e,a){return n("el-table-column",{key:a,attrs:{label:e,align:"center"},scopedSlots:t._u([{key:"header",fn:function(e){var a=e.column;return[n("span",[t._v(t._s(a.label))])]}},{key:"default",fn:function(e){var r=e.row,l=e.column;return[n("span",0===a?[t._v(t._s(Object.keys(r)[0]))]:[t._v(t._s(Object.values(r)[0][l.label]))])]}}],null,!0)})})),1)],1)},i=[],c=n("a34a"),o=n.n(c),u=n("f976");function s(t,e,n){switch(n.length){case 0:return t.call(e);case 1:return t.call(e,n[0]);case 2:return t.call(e,n[0],n[1]);case 3:return t.call(e,n[0],n[1],n[2])}return t.apply(e,n)}var f=s,h=Math.max;function b(t,e,n){return e=h(void 0===e?t.length-1:e,0),function(){var a=arguments,r=-1,l=h(a.length-e,0),i=Array(l);while(++r<l)i[r]=a[e+r];r=-1;var c=Array(e+1);while(++r<e)c[r]=a[r];return c[e]=n(i),f(t,this,c)}}var p=b;function v(t){return function(){return t}}var d=v,m=n("0305"),g=m["a"]?function(t,e){return Object(m["a"])(t,"toString",{configurable:!0,enumerable:!1,value:d(e),writable:!0})}:u["a"],y=g,w=800,C=16,S=Date.now;function _(t){var e=0,n=0;return function(){var a=S(),r=C-(a-n);if(n=a,r>0){if(++e>=w)return arguments[0]}else e=0;return t.apply(void 0,arguments)}}var j=_,D=j(y),x=D;function k(t,e){return x(p(t,e,u["a"]),t+"")}var O=k,z=n("b703"),$=n("0148");function P(t,e,n,a){var r=n-1,l=t.length;while(++r<l)if(a(t[r],e))return r;return-1}var A=P,E=n("a2fb"),T=n("7804"),I=Array.prototype,F=I.splice;function Q(t,e,n,a){var r=a?A:$["a"],l=-1,i=e.length,c=t;t===e&&(e=Object(T["a"])(e)),n&&(c=Object(z["a"])(t,Object(E["a"])(n)));while(++l<i){var o=0,u=e[l],s=n?n(u):u;while((o=r(c,s,o,a))>-1)c!==t&&F.call(c,o,1),F.call(t,o,1)}return t}var M=Q;function N(t,e){return t&&t.length&&e&&e.length?M(t,e):t}var R=N,q=O(R),H=q,J=n("8cbb");function U(t){return t&&t.length?Object(J["a"])(t):[]}var Z=U,L=n("c949");function B(t,e,n){return e in t?Object.defineProperty(t,e,{value:n,enumerable:!0,configurable:!0,writable:!0}):t[e]=n,t}function G(t,e,n,a,r,l,i){try{var c=t[l](i),o=c.value}catch(u){return void n(u)}c.done?e(o):Promise.resolve(o).then(a,r)}function K(t){return function(){var e=this,n=arguments;return new Promise((function(a,r){var l=t.apply(e,n);function i(t){G(l,a,r,i,c,"next",t)}function c(t){G(l,a,r,i,c,"throw",t)}i(void 0)}))}}var V={data:function(){return{cols:[""],keys:[""],tableData:[],timeFilter:null,refresh:null}},mounted:function(){this.fetchData()},methods:{fetchData:function(){var t=this;return K(o.a.mark((function e(){var n,a;return o.a.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return e.next=2,L["g"].zkStatus(t.$route.params.id);case 2:n=e.sent,a=n.data.entity,t.cols=[""],t.keys=[""],t.tableData=[],a.forEach((function(e){t.cols.push(e.host),t.keys=H(Object.keys(e),"host")})),t.keys.forEach((function(e){var n=B({},e,{});a.forEach((function(a){n[e][a["host"]]=a[e],t.tableData.push(n)})),t.tableData=Z(t.tableData)}));case 9:case"end":return e.stop()}}),e)})))()},timeFilterChange:function(){this.fetchData()},timeFilterRefresh:function(){this.fetchData()}}},W=V,X=n("2877"),Y=Object(X["a"])(W,l,i,!1,null,"3885750b",null),tt=Y.exports,et=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"replication-status pb-20"},[n("div",{staticClass:"title flex flex-between flex-vcenter ptb-10"},[n("span",{staticClass:"fs-20 font-bold"},[t._v(t._s(t.$t("tables.Table Replication Status")))])]),n("el-table",{staticClass:"tb-edit",staticStyle:{width:"100%"},attrs:{data:t.tableData.slice((t.currentPage-1)*t.pageSize,t.currentPage*t.pageSize),"header-cell-style":t.mergeTableHeader,border:""}},t._l(t.cols,(function(e,a){return n("el-table-column",{key:a,ref:"tableColumn",refInFor:!0,attrs:{label:e.label,prop:e.prop,width:"auto",align:"center"},scopedSlots:t._u([{key:"header",fn:function(e){var a=e.column;return[n("span",[t._v(t._s(a.label))])]}},{key:"default",fn:function(e){var r=e.row,l=e.column;return[n("span",0===a?[t._v(t._s("Table Name"===Object.keys(r)[0]?t.$t("common."+Object.keys(r)[0]):Object.keys(r)[0]))]:[t._v(t._s(Object.values(r)[0][l.property]))])]}}],null,!0)})})),1),n("div",{staticClass:"text-center"},[t.tableData.length>0?n("el-pagination",{attrs:{"current-page":t.currentPage,"page-sizes":[5,10,20,40],"page-size":t.pageSize,layout:"sizes, prev, pager, next, jumper",total:t.tableData.length},on:{"size-change":t.handleSizeChange,"current-change":t.handleCurrentChange}}):t._e()],1)],1)},nt=[],at=n("5c8a");function rt(t,e){return e="function"==typeof e?e:void 0,t&&t.length?Object(J["a"])(t,void 0,e):[]}var lt=rt,it=n("12a1");function ct(t,e){return Object(it["a"])(t,e)}var ot=ct;function ut(t,e){return pt(t)||bt(t,e)||ft(t,e)||st()}function st(){throw new TypeError("Invalid attempt to destructure non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method.")}function ft(t,e){if(t){if("string"===typeof t)return ht(t,e);var n=Object.prototype.toString.call(t).slice(8,-1);return"Object"===n&&t.constructor&&(n=t.constructor.name),"Map"===n||"Set"===n?Array.from(t):"Arguments"===n||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?ht(t,e):void 0}}function ht(t,e){(null==e||e>t.length)&&(e=t.length);for(var n=0,a=new Array(e);n<e;n++)a[n]=t[n];return a}function bt(t,e){var n=t&&("undefined"!==typeof Symbol&&t[Symbol.iterator]||t["@@iterator"]);if(null!=n){var a,r,l=[],i=!0,c=!1;try{for(n=n.call(t);!(i=(a=n.next()).done);i=!0)if(l.push(a.value),e&&l.length===e)break}catch(o){c=!0,r=o}finally{try{i||null==n["return"]||n["return"]()}finally{if(c)throw r}}return l}}function pt(t){if(Array.isArray(t))return t}function vt(t,e,n){return e in t?Object.defineProperty(t,e,{value:n,enumerable:!0,configurable:!0,writable:!0}):t[e]=n,t}function dt(t,e,n,a,r,l,i){try{var c=t[l](i),o=c.value}catch(u){return void n(u)}c.done?e(o):Promise.resolve(o).then(a,r)}function mt(t){return function(){var e=this,n=arguments;return new Promise((function(a,r){var l=t.apply(e,n);function i(t){dt(l,a,r,i,c,"next",t)}function c(t){dt(l,a,r,i,c,"throw",t)}i(void 0)}))}}var gt={data:function(){return{cols:[],tableData:[],headerData:[],timeFilter:null,refresh:null,currentPage:1,pageSize:10}},mounted:function(){this.fetchData()},methods:{fetchData:function(){var t=this;return mt(o.a.mark((function e(){var n,a,r,l,i,c,u;return o.a.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return e.next=2,L["g"].replicationStatus(t.$route.params.id);case 2:n=e.sent,a=n.data.entity,r=a.header,l=void 0===r?[]:r,i=a.tables,c=void 0===i?[]:i,t.cols=[{prop:"",label:""}],t.headerData=Object(at["a"])(l),t.tableData=[],u={},l.forEach((function(e,n){var a="shard".concat(n+1);e.forEach((function(e,n){u["".concat(a,"_").concat(n)]=e,t.cols.push({prop:"".concat(a,"_").concat(n),label:a})}))})),t.tableData.push(vt({},"Table Name",u)),c.forEach((function(e){var n=e.name,a=e.values,r={};a.forEach((function(t,e){var n="shard".concat(e+1);t.forEach((function(t,e){r["".concat(n,"_").concat(e)]=t}))})),t.tableData.push(vt({},n,r)),t.tableData=lt(t.tableData,ot)}));case 15:case"end":return e.stop()}}),e)})))()},mergeTableHeader:function(t){t.row,t.column;var e=t.rowIndex,n=t.columnIndex,a=new Set(this.headerData.map((function(t){return t.length}))),r=ut(a,1),l=r[0];if(0===e&&0!=n){if(n%l===0)return{display:"none"};this.$nextTick((function(){var t=document.querySelector(".replication-status thead>tr").children;t[n]&&(t[n].colSpan=2)}))}},timeFilterChange:function(){this.fetchData()},timeFilterRefresh:function(){this.fetchData()},handleSizeChange:function(t){this.pageSize=t},handleCurrentChange:function(t){this.currentPage=t}}},yt=gt,wt=(n("0a8e"),Object(X["a"])(yt,et,nt,!1,null,"5d203a77",null)),Ct=wt.exports,St=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"table-metric pb-20"},[n("div",{staticClass:"title flex flex-between flex-vcenter ptb-10"},[n("span",{staticClass:"fs-20 font-bold"},[t._v(t._s(t.$t("tables.Table Metrics")))])]),n("el-table",{attrs:{data:t.tableData.slice((t.currentPage-1)*t.pageSize,t.currentPage*t.pageSize),center:"",border:""}},[t._l(t.columns,(function(t){var e=t.prop,a=t.label;return[n("el-table-column",{key:e,attrs:{prop:e,label:a,"show-overflow-tooltip":""}})]}))],2),n("div",{staticClass:"text-center"},[t.tableData.length>0?n("el-pagination",{attrs:{"current-page":t.currentPage,"page-sizes":[5,10,20,40],"page-size":t.pageSize,layout:"sizes, prev, pager, next, jumper",total:t.tableData.length},on:{"size-change":t.handleSizeChange,"current-change":t.handleCurrentChange}}):t._e()],1)],1)},_t=[];function jt(t,e){return zt(t)||Ot(t,e)||xt(t,e)||Dt()}function Dt(){throw new TypeError("Invalid attempt to destructure non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method.")}function xt(t,e){if(t){if("string"===typeof t)return kt(t,e);var n=Object.prototype.toString.call(t).slice(8,-1);return"Object"===n&&t.constructor&&(n=t.constructor.name),"Map"===n||"Set"===n?Array.from(t):"Arguments"===n||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?kt(t,e):void 0}}function kt(t,e){(null==e||e>t.length)&&(e=t.length);for(var n=0,a=new Array(e);n<e;n++)a[n]=t[n];return a}function Ot(t,e){var n=t&&("undefined"!==typeof Symbol&&t[Symbol.iterator]||t["@@iterator"]);if(null!=n){var a,r,l=[],i=!0,c=!1;try{for(n=n.call(t);!(i=(a=n.next()).done);i=!0)if(l.push(a.value),e&&l.length===e)break}catch(o){c=!0,r=o}finally{try{i||null==n["return"]||n["return"]()}finally{if(c)throw r}}return l}}function zt(t){if(Array.isArray(t))return t}function $t(t,e,n,a,r,l,i){try{var c=t[l](i),o=c.value}catch(u){return void n(u)}c.done?e(o):Promise.resolve(o).then(a,r)}function Pt(t){return function(){var e=this,n=arguments;return new Promise((function(a,r){var l=t.apply(e,n);function i(t){$t(l,a,r,i,c,"next",t)}function c(t){$t(l,a,r,i,c,"throw",t)}i(void 0)}))}}var At={data:function(){return{tableData:[],currentPage:1,pageSize:10}},mounted:function(){this.fetchData()},methods:{fetchData:function(){var t=this;return Pt(o.a.mark((function e(){var n,a;return o.a.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return e.next=2,L["g"].tableMetrics(t.$route.params.id);case 2:n=e.sent,a=n.data.entity,Object.entries(a).forEach((function(e){var n=jt(e,2),a=n[0],r=n[1],l=r.columns,i=r.rows,c=r.space,o=r.completedQueries,u=r.failedQueries,s=r.parts,f=r.queryCost;t.tableData.push({tableName:a,columns:l,rows:i,space:c,completedQueries:o,failedQueries:u,parts:s,queryCost:Object.values(f).join(", ")})}));case 5:case"end":return e.stop()}}),e)})))()},handleSizeChange:function(t){this.pageSize=t},handleCurrentChange:function(t){this.currentPage=t}},computed:{columns:function(){var t=[{prop:"tableName",label:this.$t("tables.Table Name")},{prop:"columns",label:this.$t("tables.Columns")},{prop:"rows",label:this.$t("tables.Rows")},{prop:"parts",label:this.$t("tables.Parts")},{prop:"space",label:this.$t("tables.Disk Space")},{prop:"completedQueries",label:this.$t("tables.Completed Queries in last 24h")},{prop:"failedQueries",label:this.$t("tables.Failed Queries in last 24h")},{prop:"queryCost",label:this.$t("tables.Last 7 days info")}];return t}}},Et=At,Tt=(n("3249"),Object(X["a"])(Et,St,_t,!1,null,"0ea6cb38",null)),It=Tt.exports,Ft={data:function(){return{}},mounted:function(){},methods:{},components:{ZkTable:tt,ReplicationTable:Ct,TableMetric:It}},Qt=Ft,Mt=Object(X["a"])(Qt,a,r,!1,null,"11456861",null);e["default"]=Mt.exports},a3cf:function(t,e,n){}}]);