<?php

    $url = 'http://localhost:8060/';
    $data['servicePath'] = 'v1';
    $data['serviceMethod'] = 'queryInt';
    $data['content'] = '{"Id": 10000, "msg": "测试api"}';

    $this->curl = new CurlRequest($url);
    $this->curl->setPostFeilds($data);
    $data = $this->curl->post();
    $data = json_decode($data, true);
    var_dump($data);

    // curl
    class CurlRequest
    {
        public $version = "1.0";

        public $handler;

        public $timeOut = 30;

        public $header = ['Expect:'];

        public $referer;

        public $postFields;

        public $url;

        public $ssl = false;

        public $type = "get";

        public $agent = "see curl/1.0";

        public $returnData;

        public function __construct($url)
        {
            $this->url = $url = trim($url);
            $this->handler = curl_init();
            curl_setopt($this->handler, CURLOPT_HTTP_VERSION, CURL_HTTP_VERSION_1_0);
            curl_setopt($this->handler, CURLOPT_RETURNTRANSFER, true);
            curl_setopt($this->handler, CURLOPT_URL, $url);
            //支持https
            $this->ssl = stripos($url, 'https://') === 0 ? true : false;
            if ($this->ssl) {
                curl_setopt($this->handler, CURLOPT_SSL_VERIFYPEER, false);
                curl_setopt($this->handler, CURLOPT_SSL_VERIFYHOST, false);
            }
        }

        public function setTimeOut($timeOut){
            $this->timeOut = $timeOut;
        }

        public function setOpt($opt, $value)
        {
            curl_setopt($this->handler, $opt, $value);
        }

        public function setHeader($header)
        {
            $this->header = array_merge($this->header, $header);
        }

        public function setCookie($cookie)
        {
            curl_setopt($this->handler, CURLOPT_COOKIE, $cookie);
        }

        public function setReferer($referer)
        {
            $this->referer = $referer;
        }

        public function setPostFeilds($postFields)
        {
            $this->postFields = $postFields;
        }

        public function setAgent($agent){
            $this->agent = $agent;
        }

        public function exec()
        {
            if (!empty($this->referer)) {
                curl_setopt($this->handler, CURLOPT_REFERER, $this->referer);
            }else{
                curl_setopt($this->handler, CURLOPT_AUTOREFERER, true);
            }

            //set header
            $this->setHeader(['PHP-SEE-TID:' . 1]);
            $this->setHeader(['PHP-SEE-SEQ:' . 2]);

            curl_setopt($this->handler,CURLOPT_HTTPHEADER,$this->header);
            curl_setopt($this->handler, CURLOPT_USERAGENT,$this->agent);
            curl_setopt($this->handler, CURLOPT_HTTP_VERSION, CURL_HTTP_VERSION_1_0);
            curl_setopt($this->handler, CURLOPT_IPRESOLVE, CURL_IPRESOLVE_V4);
            curl_setopt($this->handler, CURLOPT_CONNECTTIMEOUT, $this->timeOut);
            curl_setopt($this->handler, CURLOPT_TIMEOUT, $this->timeOut);

            if($this->type=='post'){
                curl_setopt($this->handler, CURLOPT_POST, true);
                curl_setopt($this->handler, CURLOPT_POSTFIELDS, $this->postFields);
            }

            $this->returnData = curl_exec($this->handler);
            $httpcode = curl_getinfo($this->handler, CURLINFO_HTTP_CODE);
            if ($errorNo = curl_errno($this->handler) || $httpcode != 200) {
                //error message
                $errorMsg = curl_error($this->handler);
                \See::$log->warning("curl error, url:%s, type:%s, postData:%s, errorNo:%s, errorMsg:%s, httpcode:%s, return:%s",$this->url,$this->type,$this->postFields,$errorNo,$errorMsg,$httpcode,$this->returnData);
            }else{
                \See::$log->trace("curl success, url:%s, type:%s, postData:%s, httpcode:%s, return:%s",$this->url,$this->type,$this->postFields,$httpcode,$this->returnData);
            }
            curl_close($this->handler);
            return $this->returnData;
        }

        public function get(){
            $this->type = 'get';
            return $this->exec();
        }

        public function post(){
            $this->type ='post';
            return $this->exec();
        }

    }