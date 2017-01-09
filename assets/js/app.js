(function () {
    'use strict';

    var baseURL = 'http://localhost:9090/db/timesheet';

    Vue.component('SearchItems', {
        template: `
        <div class="card card-block">
            <h4 class="card-title">{{ header }} <button type="button" class="btn btn-link">+</button></h4>
            <form>
                <div class="form-group">
                    <input class="form-control form-control-sm" type="text" v-model="searchBy" placeholder="search">
                </div>
                <div class="form-check" v-if="showOwned">
                    <label class="form-check-label">
                        <input class="form-check-input" type="checkbox" v-model="owned">
                        owned
                    </label>
                </div>
            </form>
        </div>
        `,
        data: function () {
            return { searchBy: '', owned: true }
        },
        props: {
            showOwned: {
                type: Boolean,
                default: false
            },
            header: {
                type: String,
                required: true
            }
        }
    });

    Vue.component('ItemsList', {
        template: `
        <div class="card">
            <ul class="list-group list-group-flush" v-for="item in items">
                <li class="list-group-item bg-faded">{{ item.name }}</li>
                <li class="list-group-item">
                    <ul v-for="subitem in item.subitems.slice(0, 3)">
                        <li>{{ subitem.name }}</li>
                    </ul>
                    <ul v-if="item.subitems.length > 3">
                        <li>
                            <button type="button" class="btn btn-link">...</button>
                        </li>
                    </ul>
                </li>
            </ul>
        </div>
        `,
        props: ['items']
    });

    var app = new Vue({
        el: '#app',
        data: function () {
            var data = {};
            $.getJSON(baseURL + '/timesheet/user_id/1.json')
                .then(function (data) {
                    data = data;
                    console.log(data);
                })
            return data;
        }
    });
})();