<?php
require_once __DIR__ . '/vendor/autoload.php';

use Carbon\Carbon;
use Carbon\CarbonInterface;

// Pre-create objects outside the loop — matches Go benchmarks which
// create time.Time objects before the benchmark loop.
$now = Carbon::create(2025, 3, 29, 12, 0, 0, 'UTC');
$from = Carbon::create(2024, 6, 15, 10, 30, 0, 'UTC');
$from3Parts = Carbon::create(2024, 1, 15, 10, 30, 0, 'UTC');

$cases = [
    'DiffForHumans'      => fn() => $from->diffForHumans($now),
    'DiffForHumans3Parts' => fn() => $from3Parts
        ->diffForHumans($now, CarbonInterface::DIFF_RELATIVE_TO_OTHER, false, 3),
    'DiffForHumansShort' => fn() => $from
        ->diffForHumans($now, CarbonInterface::DIFF_RELATIVE_TO_OTHER, true),
];

$iterations = 50000;

foreach ($cases as $name => $fn) {
    for ($i = 0; $i < 500; $i++) $fn();

    $start = hrtime(true);
    for ($i = 0; $i < $iterations; $i++) {
        $fn();
    }
    $elapsed = hrtime(true) - $start;

    $nsPerOp = $elapsed / $iterations;
    printf("%-25s %10.1f ns/op\n", $name, $nsPerOp);
}
