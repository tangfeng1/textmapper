<?xml version="1.0" encoding="windows-1251" ?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
<xsl:output method="html" encoding="windows-1251"/>
<xsl:param name="lang"/>

<xsl:template match="hxml">
	<xsl:call-template name="main_layout"/>
</xsl:template>

<!-- ################################################################################### -->

<xsl:template name="main_layout">
<xsl:text disable-output-escaping="yes"><![CDATA[<!-- ERROR: Zero-length file! -->



































































]]></xsl:text><html>
<head>
    <title>Lexical Analyzer and Parser Generator</title>
    <link href="styles.css" type="text/css" rel="stylesheet" />
    <meta http-equiv="Content-Type" content="text/html; charset=windows-1251" />
</head>
<body bgcolor="#FFFFFF" link="#2971C1" vlink="#2971C1" alink="#2971C1" marginwidth="0"
    marginheight="0" leftmargin="0" rightmargin="0" topmargin="0" bottommargin="0">
<a name="top"></a>

<table width="100%" height="100%" cellspacing="0" cellpadding="0" border="0">
<tr><td height="95%" valign="top">

	<table width="965" cellspacing="0" cellpadding="0" border="0">
	<tr>
        <td width="40"><img src="images/empty.gif" height="1" width="40" border="0" alt=""/></td>
        <td width="100%"><img src="images/empty.gif" height="1" width="900" border="0" alt=""/></td>
        <td width="25"><img src="images/empty.gif" height="1" width="25" border="0" alt=""/></td>
	</tr>
	<tr>
    	<td></td>
    	<td>
            <table width="100%" cellspacing="0" cellpadding="0" border="0">
            <tr>
            <td width="70%" align="right" valign="top">
                <table cellspacing="0" cellpadding="0" border="0" width="100%">
                <tr>
                	<td height="23" colspan="3"></td>
                </tr>
                <tr>
                    <td align="right" valign="top">

                        <table cellspacing="0" cellpadding="0" border="0"><tr>
                        	<td valign="bottom" background="images/underlbg.gif" style="background-position:bottom;background-repeat:repeat-x;">
	                        
	                        <table cellspacing="0" cellpadding="0" border="0">
                            <tr>
<xsl:for-each select="i18n/menu[@lang=$lang]/item[position() &lt; 4]">
	<xsl:if test="position() &gt; 1">
                                <td class="textform">
                                    <img src="images/topdiv.gif" height="15" width="25" border="0" alt="|"/></td>
	</xsl:if>
                                <td class="textform">
                                    <a><xsl:attribute name="href"><xsl:value-of select="@link"/></xsl:attribute><xsl:attribute name="title"><xsl:value-of select="@title"/></xsl:attribute><xsl:value-of select="@title"/></a></td>
</xsl:for-each>
                            </tr>
            	            <tr>
            	            	<td height="5" colspan="5"></td>
            	            </tr>
	                        </table>
                        </td></tr></table>
                    </td>
                    <td width="20">
                        <img src="images/empty.gif" height="1" width="20" border="0" alt=""/></td>
                	<td width="155" align="right">
                		<img src="images/head_left.jpg" height="60" width="155" border="0" alt=""/></td>
                </tr>
                <tr>
                    <td colspan="3" height="54" valign="bottom" align="left">

                        <table cellspacing="0" cellpadding="0" border="0">
                        <tr>
	    	                <td width="6">
    	    	                <img src="images/empty.gif" height="1" width="6" border="0" alt=""/></td>
                        	<td>
								<span class="head1"><xsl:value-of select="title[@lang=$lang]"/></span>
                        	</td>
                        </tr>
                        <tr>
                        	<td></td>
                        	<td valign="bottom" height="7">
                    			<img src="images/red_lineh.gif" height="3" width="243" border="0" alt=""/></td>
						</tr>
                        <tr>
                        	<td></td>
                        	<td height="3"></td>
                        </tr>
						</table>
					</td>
                </tr>
                </table>
            </td>
            <td align="center" valign="top">
                <table cellspacing="0" cellpadding="0" border="0">
                <tr>
                	<td height="17"></td>
                </tr>
                <tr>
                	<td>
		                <img src="images/red_line.gif" height="120" width="5" border="0" alt=""/></td>
		        </tr>
		        </table>
            </td>
            <td width="450" align="left" valign="top">
                <table cellspacing="0" cellpadding="0" border="0">
                <tr>
                    <td align="right" height="45">
                        <table cellspacing="0" cellpadding="0" border="0">
                        	<tr>
                        		<td height="5"></td>
                        	</tr>
                            <tr>
                                <td><a href="index.html" title="�� �������">
                                    <img src="images/icon_home.gif" height="15" width="15" border="0" alt="�� �������"/></a></td>
                                <td>
                                    <img src="images/icondiv.gif" height="30" width="55" border="0" alt="|"/></td>
                                <td><a href="mailto:eugeniy@gryaznov.net">
                                    <img src="images/icon_mail.gif" height="10" width="15" border="0" alt="��������� ������"/></a></td>
                                <td width="10" height="40"><img src="images/empty.gif" height="40" width="10" border="0" alt=""/></td>
                            </tr>
                        </table>
                    </td>
                </tr>
                <tr>
                	<td><img src="images/head_right.jpg" height="60" width="429" border="0" alt=""/></td>
                </tr>
                <tr>
					<td>
                        <table cellspacing="0" cellpadding="0" border="0">
                        <tr>
				            <td width="30" height="32">
                		        <img src="images/empty.gif" height="32" width="30" border="0" alt=""/></td>
                		    <td valign="bottom" align="left" background="images/underlbg.gif" 
                		    	style="background-position:bottom;background-repeat:repeat-x;">
                		        <table cellspacing="0" cellpadding="0" border="0">
                        		<tr>
<xsl:for-each select="i18n/menu[@lang=$lang]/item[position() &gt; 3]">
	<xsl:if test="position() &gt; 1">
                                	<td class="textform">
                                    	<img src="images/topdiv.gif" height="15" width="25" border="0" alt="|"/></td>
	</xsl:if>
                                	<td class="textform">
                                    	<a><xsl:attribute name="href"><xsl:value-of select="@link"/></xsl:attribute><xsl:attribute name="title"><xsl:value-of select="@title"/></xsl:attribute><xsl:value-of select="@title"/></a></td>
</xsl:for-each>
            	                </tr>
            	                <tr>
            	                	<td height="5" colspan="5"></td>
            	                </tr>
			                    </table>
                            </td>
                        </tr>
                        </table>
					</td>
                </tr>
                </table>
            </td>
            </tr>
            </table>

    	</td>
        <td></td>
	</tr>
	<tr>
        <td></td>
        <td>
            <table width="100%" cellspacing="0" cellpadding="0" border="0">
            <tr>
            	<td valign="top">
            	    <table width="100%" cellspacing="0" cellpadding="0" border="0">
            	    <tr>
            	    	<td height="18"></td>
            	    </tr>
            	    <tr>
            	    	<td class="text">
            	    		<xsl:apply-templates select="content"/>
            	    	</td>
            	    </tr>
            	    </table>
            	</td>
            </tr>
            </table>
        
        
        </td>
        <td></td>
	</tr>
	</table>
	
</td></tr><tr><td>

    <table width="100%" cellspacing="0" cellpadding="0" border="0"><tr><td valign="bottom" background="images/grey.gif" style="background-position:bottom;background-repeat:repeat-x;">
        <table width="965" cellspacing="0" cellpadding="0" border="0">
        <tr>
	        <td width="40"><img src="images/empty.gif" height="1" width="40" border="0" alt=""/></td>
    	    <td width="100%"><img src="images/empty.gif" height="1" width="900" border="0" alt=""/></td>
        	<td width="25"><img src="images/empty.gif" height="1" width="25" border="0" alt=""/></td>
        </tr>
        <tr>
            <td height="40"></td>
            <td></td>
            <td></td>
        </tr>
        <tr>
            <td></td>
            <td valign="bottom">
                <table width="100%" cellspacing="0" cellpadding="0" border="0"><tr>
                <td width="25%">
	                <table cellspacing="0" cellpadding="0" border="0"><tr>
    	                <td>
        	                <img src="images/icon_address.gif" height="10" width="20" border="0" alt=""/></td>
                	    <td class="address">
                    	    <a href="mailto:eugeniy@gryaznov.net" title="eugeniy@gryaznov.net">eugeniy<i>@</i>gryaznov.net</a></td>
	                </tr>
    	            </table>
				</td>
                <td width="25%" valign="bottom" align="right">
                    <table cellspacing="0" cellpadding="0" border="0"><tr>
                    <td>
                        <a href="#top" title="������ ��������">
                            <img src="images/butt_top.gif" height="40" width="60" border="0" alt=""/></a></td>
                    <td class="address">
                        <a href="#top" title="������ ��������">������</a></td>
                    </tr>
                    </table>
                </td>
                <td width="25%" valign="bottom">
                </td>
                <td width="25%" valign="bottom" class="phone" align="right">
                    <img src="images/empty.gif" height="3" width="1" border="0" alt=""/><br/>
                    <span class="number">+7 (921) 3497513</span><br/>
                    <img src="images/empty.gif" height="10" width="1" border="0" alt=""/><br/>
                </td>
                </tr>
                </table>
            </td>
            <td></td>
        </tr>
        </table>
    </td></tr>
    </table>
</td></tr>
</table>

</body>
</html>


</xsl:template>

<!-- ################################################################################### -->

<xsl:template match="text">
	<xsl:value-of select="text()" disable-output-escaping="yes"/>
</xsl:template>

<xsl:template match="h">
	<span class="head3"><xsl:apply-templates/></span>
</xsl:template>

<xsl:template match="i|b|p|br|li|ul">
	<xsl:copy>
		<xsl:apply-templates/>
	</xsl:copy>
</xsl:template>

<xsl:template match="a">
	<xsl:copy><xsl:copy-of select="@*"/>
		<xsl:apply-templates/>
	</xsl:copy>
</xsl:template>

<xsl:template match="lang">
	<xsl:if test="$lang=@lang">
	<xsl:apply-templates/>
	</xsl:if>
</xsl:template>

<!-- ################################################################################### -->

<!-- DOWNLOAD -->

<xsl:template match="download">
<xsl:variable name="dwlddir" select="@dir"/>

<span class="head3">Download</span>

<p>
<xsl:for-each select="group">
<table cellspacing="0" cellpadding="0" border="0">
<tr>
<td width="20"><img src="images/empty.gif" height="1" width="20" border="0" alt=""/></td>
<td width="100"><img src="images/empty.gif" height="1" width="100" border="0" alt=""/></td>
<td width="100%"><img src="images/empty.gif" height="1" width="200" border="0" alt=""/></td>
</tr>
<tr><td colspan="3" class="loadh"><xsl:value-of select="@name"/></td></tr>
<tr><td colspan="3" height="5"></td></tr>
<xsl:for-each select="file">
<tr>
<td align="center"><img src="images/dir-file.png" border="0"/></td>
<td><a><xsl:attribute name="href"><xsl:value-of select="$dwlddir"/><xsl:value-of select="@name"/></xsl:attribute><xsl:value-of select="@name"/></a></td>
<td><xsl:value-of select="."/></td>
</tr>
</xsl:for-each>
</table>
</xsl:for-each>
</p>

</xsl:template>

<!-- LINKS -->

<xsl:template match="links">
<xsl:if test="link">
<span class="head3">Links</span>

<p>
<xsl:for-each select="link">
<a><xsl:attribute name="href"><xsl:value-of select="@href"/></xsl:attribute><xsl:value-of select="."/></a><br/>
</xsl:for-each>
</p>
</xsl:if>
</xsl:template>

</xsl:stylesheet>
