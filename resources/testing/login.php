<DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<!-- Google Tag Manager -->
<script>(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start':
new Date().getTime(),event:'gtm.js'});var f=d.getElementsByTagName(s)[0],
j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src=
'https://www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);
})(window,document,'script','dataLayer','GTM-M867HWJ');</script>
<!-- End Google Tag Manager -->
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta http-equiv="Content-Language" content="en-usa"/>
<meta name="robots" content="noindex,nofollow"/>
<link rel="stylesheet" type="text/css" href="//instantproxies.com/billing/templates/default/css/bootstrap.css" />
<link rel="stylesheet" type="text/css" href="//instantproxies.com/billing/templates/default/css/whmcs.css" />
<link rel="shortcut icon" href="img/favicon.png"/>
<script type="text/javascript" src="//instantproxies.com/monitor/stats.js"></script>
<script type="text/javascript" src="//instantproxies.com/js/common.js"></script>
<title>InstantProxies.com: Proxy Admin Panel</title>
<style>
#warning_msg {
	margin-left:auto;
	margin-right:auto;
	margin-top:15px;
	margin-bottom:15px;
	width:450px;
	padding:10px;
	overflow:hidden;
	background: #FFFFE0;
	border:solid 1px #E6DB55;
	-moz-border-radius:5px;
	-webkit-border-radius: 5px;
	box-sizing: border-box;
	line-height: 166%;
}
.reset_input {
	width: 200px;
	display:block;
	float:left;
	margin-bottom:12px;
	padding: 2px;
	border:solid 1px #d3d3d3;
	-moz-border-radius: 3px;
	-webkit-border-radius: 3px;
	border-radius: 3px;
}
label {
  display: block;
  width: 80px;
  float: left;
  padding-top: 10px;
}
</style>
<script>
function initReset() {
	document.getElementById('reset_but').style.display = 'none';
	document.getElementById('reset_section').style.display = '';
}
</script>
</head>

<body>
<div id="top_container">
  <div id="header">
  <div id="header_iner">
    <div id="logo" style="padding:12px 0 0"><a href="//instantproxies.com/"><img src="//instantproxies.com/images/logo2.png" alt="logo" width="350" height="68" border="0" /></a></div>

    <div id="header_right">
      <div id="header_right_1" style="width:635px;">
        <div align="right" class="top_link"><a href="//instantproxies.com/billing/clientarea.php"><img src="//instantproxies.com/billing/templates/default/images/lock.jpg" alt="lock" width="14" height="18" border="0" style="float:left; margin:0 5px;" />Client Login</a>&nbsp;&nbsp;&nbsp;&nbsp;<a href="javascript:openSupportChat()"><img src="//instantproxies.com/billing/templates/default/images/chat.jpg" alt="chat" width="20" height="18" border="0" /> Live Support</a></div>
      </div>
      <!--
      <div id="header_right_2">
        <div align="right"><span class="ip_nav"><ul><li><a href="//instantproxies.com/">Home</a></li>
        <li><a href="//instantproxies.com/pricing/">Pricing</a></li>

        <li class="{php} echo $activeTab1{/php}"><a href="//instantproxies.com/billing/knowledgebase.php">FAQ</a></li>
        <li class="navr"><a href="//instantproxies.com/resources/">Resources</a></li>
        <li><a href="//instantproxies.com/affiliates/">Affiliates</a></li>
        <li class="{php} echo $activeTab2{/php}"><a href="//instantproxies.com/billing/submitticket.php?step=2&deptid=2">Support</a></li>
        </ul></span></div>
      </div>
      -->
    </div>

  </div>
</div>

<div id="banner_iner_page">
  <div id="banner_iner_2">
    <div class="banner_heading" id="banner_iner_text" style="clear:both; padding-left: 65px">
      <h1 class="entry-title">Proxy Admin Panel</h1>
    </div>
  </div>
</div>
<div class="whmcscontainer">
<div class="contentpadded">

  <h3>Login to Manage Your Proxies</h3>
  To manage your billing profile or pay invoices, please login to the <a href="https://www.instantproxies.com/billing/clientarea.php?action=products">billing area</a> instead.<br>
  <br>


	<div style="margin: 20px auto; width: 300px;">
    <form id="login" method="post" action="login_do.php">
    <fieldset>
    <label for="user">User ID:</label><input type="text" id="user" name="username" value="" style="width: 200px" /><br />
    <label for="password">Password:</label><input type="password" id="password" name="password" value="" style="width: 200px" /><br />
    <input id="x" type="submit" name="button" value="Sign In" class="btn btn-primary btn-large" style="margin-left: 90px" />
    </fieldset>
    </form>
  </div>

  <br>
  <center>
  <b>Forgot your proxy panel username or password?</b><br><a href="https://instantproxies.com/billing/clientarea.php?action=products">Login to the billing area</a> to view your credentials.
  </center>

</div></div>
<div class="footerdivider">
  <div class="fill"></div>
</div>

<div style="height:20px;">&nbsp; </div>
<div id="footer" style="text-align:left;">
  <div id="footer_iner">
    <div class="nav_heading" id="footer_nav_box_1">
      <span class="nav_foter_heading">SERVICES<br></span>
      <span class="nav_foter"><a href="//instantproxies.com/pricing">Private Proxies</a><br />
        <a href="//instantproxies.com/affiliates">Affiliate Program</a> <br />
        <a href="//instantproxies.com/affiliates#resellers">Reseller Program</a> </span>
    </div>
    <div class="nav_heading" id="footer_nav_box_1">
      <span class="nav_foter_heading">INFORMATION<br /></span>
      <span class="nav_foter"><a href="//instantproxies.com/billing/knowledgebase.php">Frequent Questions</a><br />
        <a href="//instantproxies.com/resources/">Proxy Resources</a><br />
        <a href="//instantproxies.com/check-availability">Availability Checker</a></span>
    </div>

    <div class="nav_heading" id="footer_nav_box_1">
      <span class="nav_foter_heading">CUSTOMERS<br /></span>
      <span class="nav_foter"><a href="//instantproxies.com/billing/clientarea.php">Client Login</a><br />
        <a href="//instantproxies.com/terms/">Terms of Service</a><br />
        <a href="//instantproxies.com/support">Contact Support</a></span>
    </div>
    <div class="nav_heading" id="footer_nav_box_2">
      <span class="nav_foter_heading">NETWORK STATUS</span><br />
      <span class="nav_foter">
        <script>
          document.write('<a href="//instantproxies.com/network-status#current">Current status: ' + STATS_CURRENT + ' online</a><br>');
          document.write('<a href="//instantproxies.com/network-status#pastday">Past 24 hours: ' + STATS_24H + ' uptime</a><br>');
          document.write('<a href="//instantproxies.com/network-status#pastmonth">Past 1 month: ' + STATS_1M + ' uptime</a>');
        </script>
      </span>
    </div>
    <div style="clear:both"></div>
  </div>
</div>
<div class="footerline" style="clear:both">
  <span class="footer">&copy; 2022 InstantProxies.com. All Rights Reserved.</span><br />
</div>

<script type="text/javascript">
  var _gaq = _gaq || [];
  _gaq.push(['_setAccount', 'UA-35673112-1']);
  _gaq.push(['_setDomainName', 'instantproxies.com']);
  _gaq.push(['_trackPageview']);

  (function() {
    var ga = document.createElement('script');
    ga.type = 'text/javascript';
    ga.async = true;
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    var s = document.getElementsByTagName('script')[0];
    s.parentNode.insertBefore(ga, s);
  })();
</script>
<!--Start of Tawk.to Script-->
<script type="text/javascript">
  var Tawk_API = Tawk_API || {},
    Tawk_LoadStart = new Date();
  (function() {
    var s1 = document.createElement("script"),
      s0 = document.getElementsByTagName("script")[0];
    s1.async = true;
    s1.src = 'https://embed.tawk.to/6274369fb0d10b6f3e70d13c/1g2avaovc';
    s1.charset = 'UTF-8';
    s1.setAttribute('crossorigin', '*');
    s0.parentNode.insertBefore(s1, s0);
  })();
</script>
<!--End of Tawk.to Script-->
</body>
</html>
