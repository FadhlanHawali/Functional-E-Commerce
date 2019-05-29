<template>
  <div class="container">
    <nav class="navbar" role="navigation" aria-label="main navigation">
    <div id="navbarBasicExample" class="navbar-menu">
    <div class="navbar-brand">
        <a class="navbar-item" href="#">
          <h1>Immerce</h1>
        </a>
    </div>
      <div class="navbar-start">
        <a class="navbar-item" href="documentation.html">
          Documentation
        </a>
      </div>

      <div class="navbar-end">
        <div class="navbar-item">
          <div class="buttons">
            <a class="button is-primary" href="/api.js" download>
              <strong>Download SDK</strong>
            </a>
          </div>
        </div>
      </div>
    </div>
  </nav>
  <div class="hello">
    <h1 style="font-size: 50px; font-weight: 800;">Immerce, e-commerce API Solution</h1>
    <div style="padding-top: 50px;">
    <h1>{{ msg }}</h1>
    <section v-if="stage == 'register'" style="padding-top: 20px;">
        <b-field label="Email">
            <b-input v-model="email" type="email" maxlength="30"></b-input>
        </b-field>

        <b-field label="Password">
            <b-input v-model="password" type="password" maxlength="40"></b-input>
        </b-field>
        <a class="button is-primary" @click="submit">Register Now</a>
        <h5 class="padding-top:5px;">Or <a href="#" @click="stage = 'signin'">Sign In</a> if you already have an acoount</h5>
    </section>
    <section v-if="stage == 'signin'">
        <b-field label="Email" :type="isEmailFieldDanger" :message="emailFieldMessage">
            <b-input v-model="email" type="email" maxlength="30"></b-input>
        </b-field>

        <b-field label="Password" :type="isPasswordFieldDanger" :message="passwordFieldMessage">
            <b-input v-model="password" type="password" maxlength="40"></b-input>
        </b-field>
        <a class="button is-primary" @click="signin">SignIn</a>
    </section>
    </div>
  </div>
  </div>
</template>

<script>
export default {
  name: 'LandingPage',
  computed: {
    isEmailFieldDanger() {
      if (this.emailFieldMessage != '') {
        return "is-danger"
      }
      return "is-success"
    },
    isPasswordFieldDanger() {
      if (this.passwordFieldMessage != '') {
        return "is-danger"
      }
      return "is-success"
    }
  },
  data () {
    return {
      stage: 'register',
      msg: 'Register today!',
      email: '',
      password: '',
      emailFieldMessage: '',
      passwordFieldMessage: '',
    }
  },
  methods: {
    submit() {
      axios.post('/api/v1/user/create', {
        email: this.email,
        password: this.password
      })
      .then(res => {
        window.sessionStorage.setItem('access_token', res.data.Data.ApiKey);
        this.$router.push({ path: 'dashboard' })
      })
      .catch(err => {
        console.log(err);
      })
    },
    signin() {
      axios.post('/api/v1/user/login', {
        email: this.email,
        password: this.password
      })
      .then(res => {
        window.sessionStorage.setItem('access_token', res.data.Data.ApiKey);
        this.$router.push({ path: 'dashboard' })
      })
      .catch(err => {
        this.emailFieldMessage = 'Email/Password Salah';
        this.passwordFieldMessage = 'Email/Password Salah';
      })
    }
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
.hello {
  margin:0 auto;
  width: 600px;
}
</style>
