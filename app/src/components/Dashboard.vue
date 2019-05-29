<template>
    <div class="container">
        <a class="button is-primary" @click="display = true">Create Store</a>
        <section v-if="display" style="padding: 20px;">
            <b-field label="Store Name">
                <b-input v-model="form.name"></b-input>
            </b-field>

            <b-field label="Store Address">
                <b-input v-model="form.address"></b-input>
            </b-field>

            <b-field label="Store Contact">
                <b-input v-model="form.contact"></b-input>
            </b-field>

            <b-field label="Bank Number">
                <b-input v-model="form.bank"></b-input>
            </b-field>

            <a class="button is-primary" @click="create">Create</a>
            <a class="button is-primary" @click="display = false">Cancel</a>
        </section>
        <b-table
            :data="list"
            :striped="true"
            :hoverable="true" style="padding-top: 35px;">
            <template slot-scope="props">
                <b-table-column field="id" label="ID" width="40" numeric>
                    {{ props.row.Id }}
                </b-table-column>

                <b-table-column field="store_name" label="Store Name">
                    {{ props.row.namaLapak }}
                </b-table-column>

                <b-table-column field="store_address" label="Store Address">
                    {{ props.row.alamat }}
                </b-table-column>

                <b-table-column field="rekening" label="Store Bank Number">
                    {{ props.row.nomerRekening }}
                </b-table-column>

                <b-table-column field="rekening" label="Store Contact">
                    {{ props.row.telepon }}
                </b-table-column>

                <b-table-column label="Action">
                    <a href="#" @click="changeKey(props.row.Id)">Api Key</a> <span> | </span>
                    <a href="#" @click="deleteStore(props.row.Id)">Delete</a>
                </b-table-column>
            </template>
            <template slot="empty">
                <section class="section">
                    <div class="content has-text-grey has-text-centered">
                        <p>
                            <b-icon
                                icon="emoticon-sad"
                                size="is-large">
                            </b-icon>
                        </p>
                        <p>Nothing here.</p>
                    </div>
                </section>
            </template>
        </b-table>
    </div>
</template>

<script>
export default {
    name: 'dashboard',
    data() {
        return {
            list: [],
            display: false,
            form: {
                name: '',
                address: '',
                contact: '',
                bank: ''
            }
        }
    },
    methods: {
        changeKey(id) {
            axios.get('/api/v1/store/'+id, {
                headers: {
                    "Authorization": window.sessionStorage.getItem('access_token')
                }
            })
            .then(res => {
                this.prompt(res.data.Data.ApiKey)
            })
            .catch(err => console.log(err))
        },
        deleteStore(id) {
            axios.delete('/api/v1/store/'+id, {
                headers: {
                    "Authorization": window.sessionStorage.getItem('access_token')
                }
            })
            .then(res => {
                this.fetch()
            })
            .catch(err => console.log(err))
        },
        create() {
            axios.post('/api/v1/store', {
                namaLapak: this.form.name,
                alamat: this.form.address,
                telepon: this.form.contact,
                nomerRekening: this.form.bank
            }, {
                headers: {
                    "Authorization": window.sessionStorage.getItem('access_token')
                }
            })
            .then(res => {
                this.prompt(res.data.Data.ApiKey);
                this.fetch();
            })
            .catch(err => {
                console.log(err);
            })
        },
        prompt(text) {
            this.$dialog.prompt({
                message: `Copy the api_key below`,
                inputAttrs: {
                    value: text,
                    disabled: true,
                },
                onConfirm: (value) => {this.$clipboard(value);this.$toast.open(`Copied!`)}
            })
        },
        fetch() {
            axios.get('/api/v1/store', {
                headers: {
                    "Authorization": window.sessionStorage.getItem('access_token')
                }
            })
            .then(res => {
                console.log(res.data);
                this.list = res.data.Data;
            })
            .catch(err => {
                console.log(err);
            })
        }
    },
    created() {
        this.fetch();
    },
}
</script>
<style scoped>

</style>