<?php
require_once __DIR__ . '/vendor/autoload.php';

use Illuminate\Support\Number;

$cases = [
    'Abbreviate'      => fn() => Number::abbreviate(489939, maxPrecision: 2),
    'AbbreviateLarge'  => fn() => Number::abbreviate(1e18),
    'ForHumans'        => fn() => Number::forHumans(489939, maxPrecision: 2),
    'AbbreviateSmall'  => fn() => Number::abbreviate(42),
    'AbbreviateNeg'    => fn() => Number::abbreviate(-489939, maxPrecision: 2),
];

$iterations = 100000;

foreach ($cases as $name => $fn) {
    // Warm up
    for ($i = 0; $i < 1000; $i++) $fn();

    $start = hrtime(true);
    for ($i = 0; $i < $iterations; $i++) {
        $fn();
    }
    $elapsed = hrtime(true) - $start;

    $nsPerOp = $elapsed / $iterations;
    printf("%-25s %10.1f ns/op\n", $name, $nsPerOp);
}
