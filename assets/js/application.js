require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");

$(() => {
    $(document).on("submit", "form", {}, (e) => {
        $(".progress").show();
    });
});