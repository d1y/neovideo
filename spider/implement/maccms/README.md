参考: https://github.com/magicblack/maccms10/blob/master/%E8%AF%B4%E6%98%8E%E6%96%87%E6%A1%A3/API%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E.txt

api接口仅供提供数据

视频接口同时支持老板xml格式的数据，增加参数 &at=xml即可。

1,视频部分
列表http://域名/api.php/provide/vod/?ac=list
详情http://域名/api.php/provide/vod/?ac=detail
同样支持老板xml格式的数据
列表api.php/provide/vod/at/xml/?ac=list
详情api.php/provide/vod/at/xml/?ac=detail

2,文章部分
列表http://域名/api.php/provide/art/?ac=list
详情http://域名/api.php/provide/art/?ac=detail

3,演员部分
列表http://域名/api.php/provide/actor/?ac=list
详情http://域名/api.php/provide/actor/?ac=detail

4,角色部分
列表http://域名/api.php/provide/role/?ac=list
详情http://域名/api.php/provide/role/?ac=detail

5,网址部分
列表http://域名/api.php/provide/website/?ac=list
详情http://域名/api.php/provide/website/?ac=detail

列表数据格式：

```json
{"code":1,"msg":"数据列表","page":1,"pagecount":1,"limit":"20","total":15,"list":[{"vod_id":21,"vod_name":"测试1","type_id":6,"type_name":"子类1","vod_en":"qingjian","vod_time":"2018-03-29 20:50:19","vod_remarks":"超清","vod_play_from":"youku"},{"vod_id":20,"vod_name":"测试2","type_id":6,"type_name":"子类1","vod_en":"baolijiequ","vod_time":"2018-03-27 21:17:52","vod_remarks":"超清","vod_play_from":"youku"},{"vod_id":19,"vod_name":"测试3","type_id":6,"type_name":"子类3","vod_en":"chaofanzhizhuxia2","vod_time":"2018-03-27 21:17:51","vod_remarks":"高清","vod_play_from":"youku"},{"vod_id":18,"vod_name":"测试4","type_id":6,"type_name":"子类4","vod_en":"muxingshangxing","vod_time":"2018-03-27 21:17:37","vod_remarks":"高清","vod_play_from":"youku"},{"vod_id":15,"vod_name":"测试5","type_id":6,"type_name":"子类5","vod_en":"yingxiongbense2018","vod_time":"2018-03-22 16:09:17","vod_remarks":"高清","vod_play_from":"qiyi,sinahd"},{"vod_id":13,"vod_name":"测试6","type_id":8,"type_name":"子类6","vod_en":"piaoxiangjianyu","vod_time":"2018-03-21 20:37:52","vod_remarks":"全36集","vod_play_from":"youku,qiyi"},{"vod_id":14,"vod_name":"测试7","type_id":8,"type_name":"子类7","vod_en":"guaitanzhimeiyingjinghun","vod_time":"2018-03-20 21:32:27","vod_remarks":"高清","vod_play_from":"qiyi"}]}
```


列表接收参数：
ac=list
t=类别ID
pg=页码
wd=搜索关键字
h=几小时内的数据
例如： http://域名/api.php/provide/vod/?ac=list&t=1&pg=5   分类ID为1的列表数据第5页


内容数据格式：

```
{"code":1,"msg":"数据列表","page":1,"pagecount":1,"limit":"20","total":1,"list":[{"vod_id":21,"vod_name":"测试1","type_id":6,"type_name":"子类1","vod_en":"qingjian","vod_time":"2018-03-29 20:50:19","vod_remarks":"超清","vod_play_from":"youku","vod_pic":"https:\/\/localhost\/view\/photo\/s_ratio_poster\/public\/p2259384068.jpg","vod_area":"大陆","vod_lang":"国语","vod_year":"2018","vod_serial":"0","vod_actor":"主演们","vod_director":"导演","vod_content":"这可是详情介绍啊","vod_play_url":"正片$http:\/\/localhost\/v_show\/id_XMTM0NTczNDExMg==.html"}]}
```



内容接收参数：
参数 ids=数据ID，多个ID逗号分割。
     t=类型ID
     pg=页码
     h=几小时内的数据

例如:   http://域名/api.php/provide/vod/?ac=detail&ids=123,567     获取ID为123和567的数据信息
        http://域名/api.php/provide/vod/?ac=detail&h=24     获取24小时内更新数据信息


另附上xml返回格式：
列表数据格式：

```xml
<?xml version="1.0" encoding="utf-8"?><rss version="5.0"><list page="1" pagecount="23" pagesize="20" recordcount="449"><video><last>2012-05-06 13:32:28</last><id>493</id><tid>9</tid><name><![CDATA[测试]]></name><type>子类1</type><dt>dplayer</dt><note><![CDATA[]]></note><vlink><![CDATA[http://localhost/vod/?493.html]]></vlink><plink><![CDATA[http://localhost/vodplay/?493-1-1.html]]></plink></video></list><class><ty id="1">分类1</ty><ty id="2">分类2</ty><ty id="3">分类3</ty><ty id="4">分类4</ty><ty id="5">子类1</ty><ty id="6">子类2</ty><ty id="7">子类3</ty><ty id="8">子类4</ty><ty id="9">子类5</ty><ty id="10">子类6</ty><ty id="11">子类7</ty><ty id="12">子类8</ty><ty id="13">子类9</ty><ty id="14">子类10</ty><ty id="15">子类11</ty></class></rss>

```

内容数据格式：

```xml
<?xml version="1.0" encoding="utf-8"?><rss version="5.0"><list page="1" pagecount="1" pagesize="20" recordcount="1"><video><last>2012-05-06 13:32:28</last><id>493</id><tid>9</tid><name><![CDATA[测试1]]></name><type>恐怖片</type><pic>http://localhost/uploads/20091130205750222.JPG</pic><lang>英语</lang><area>欧美</area><year>2012</year><state>0</state><note><![CDATA[]]></note><type>_9</type><actor><![CDATA[]]></actor><director><![CDATA[Ryan Schifrin]]></director><dl><dd from="qvod"><![CDATA[第1集$http://localhost/1.mp4|]]></dd></dl><des><![CDATA[<p>简单介绍。 <br /></p>]]></des><vlink><![CDATA[http://localhost/vod/?493.html]]></vlink><plink><![CDATA[http://localhost/vodplay/?493-1-1.html]]></plink></video></list></rss>
```