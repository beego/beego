<?php
session_id() or session_start();
if (!isset($_SESSION['first_time'])) {
  $_SESSION['first_time'] = time();
}
$data = array();
$HTTP = array();
foreach ($_SERVER as $head => $value) {
  if (strpos($head, "HTTP_") === 0) {
    $HTTP[$head] = $value;
  }
}
$data['HTTP'] = $HTTP;
$data['GET'] = $_GET;
$data['POST'] = $_POST;
$data['REQUEST'] = $_REQUEST;
$data['SESSION'] = $_SESSION;
$data['COOKIE'] = $_COOKIE;
$data['FILES'] = $_FILES;
$data['SERVER'] = $_SERVER;
$data['ENV'] = $_ENV;

//echo json_encode($data);
//echo "<pre>";
print_r($data);
//echo "</pre>";
