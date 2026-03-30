FROM php:8.3-cli-alpine

RUN apk add --no-cache unzip curl icu-dev icu-libs \
    && docker-php-ext-install intl

COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

WORKDIR /bench

COPY composer.json ./
RUN composer install --no-dev --optimize-autoloader --no-interaction

COPY bench_abbreviate.php bench_diffforhumans.php ./

CMD echo "=== Number::abbreviate() ===" && \
    php bench_abbreviate.php && \
    echo "" && \
    echo "=== Carbon::diffForHumans() ===" && \
    php bench_diffforhumans.php
