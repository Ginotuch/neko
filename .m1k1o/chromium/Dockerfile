FROM m1k1o/neko:base

ARG SRC_URL="https://github.com/macchrome/linchrome/releases/download/v86.0.4240.111-r800218-portable-ungoogled-Lin64/ungoogled-chromium_86.0.4240.111_1.vaapi_polly_linux.tar.xz"

#
# install custom chromium build from woolyss with support for hevc/x265
RUN set -eux; apt-get update; \
    apt-get install -y --no-install-recommends libatk1.0-0 libatk-bridge2.0-0 libatomic1 \
    libcups2 libgtk-3-0 libnss3 libpci3 libxcomposite1 libxss1 openbox xz-utils; \
    wget -O - /tmp/chromium.tar.xz "${SRC_URL}" | tar -xJf- -C /tmp; \
    mv /tmp/ungoogled-chromium_* /usr/lib/chromium; \
    #
    # make required changes for sandbox mode
    mv /usr/lib/chromium/chrome_sandbox /usr/lib/chromium/chrome-sandbox; \
    chown root:root /usr/lib/chromium/chrome-sandbox; \
    chmod 4755 /usr/lib/chromium/chrome-sandbox; \
    #
    # clean up
    apt-get --purge autoremove -y xz-utils; \
    apt-get clean -y; \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/*

#
# copy configuation files
COPY supervisord.conf /etc/neko/supervisord/chromium.conf
COPY preferences.json /usr/lib/chromium/master_preferences
COPY policies.json /etc/chromium/policies/managed/policies.json
COPY openbox.xml /etc/neko/openbox.xml

#
# copy extensions and policy files
COPY extensions /usr/share/chromium/extensions
