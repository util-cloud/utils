#!/usr/bin/perl
use strict;
use warnings;
use LWP::UserAgent;

my $ua = LWP::UserAgent->new;
my $release_center_url = "http://ipaddr:22205/v1/setInfoByApp";
my $req = HTTP::Request->new(POST => $release_center_url);
$req->header('content-type' => 'application/json');

my $fn = $ARGV[0];
open(my $f, '<:encoding(UTF-8)', $fn) or die "Could not open file '$fn' $!";
my $text = join('', <$f>);

$req->content($text);
print $req->as_string;
my $resp = $ua->request($req);

if ($resp->is_success) {
 my $message = $resp->decoded_content;
 print "Received reply: $message\n";
} else {
 print "HTTP POST error code: ", $resp->code, "\n";
 print "HTTP POST error message: ", $resp->message, "\n";
}
