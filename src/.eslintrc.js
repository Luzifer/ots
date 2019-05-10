// https://eslint.org/docs/user-guide/configuring

module.exports = {
  'root': true,
  'parserOptions': {
    parser: 'babel-eslint',
    sourceType: 'module',
  },
  'env': {
    node: true,
  },
  'extends': [
    /*
     * https://github.com/vuejs/eslint-plugin-vue#priority-a-essential-error-prevention
     * consider switching to `plugin:vue/strongly-recommended` or `plugin:vue/recommended` for stricter rules.
     */
    'plugin:vue/essential',
    // https://github.com/standard/standard/blob/master/docs/RULES-en.md
    'eslint:recommended',
  ],
  // required to lint *.vue files
  'plugins': ['vue'],
  'globals': {
    locale: true,
    process: true,
    version: true,
  },
  // add your custom rules here
  'rules': {
    'array-bracket-newline': ['error', { multiline: true }],
    'array-bracket-spacing': ['error'],
    'arrow-body-style': ['error', 'as-needed'],
    'arrow-parens': ['error', 'as-needed'],
    'arrow-spacing': ['error', { before: true, after: true }],
    'block-spacing': ['error'],
    'brace-style': ['error', '1tbs'],
    'comma-dangle': ['error', 'always-multiline'], // Apply Contentflow rules
    'comma-spacing': ['error'],
    'comma-style': ['error', 'last'],
    'curly': ['error'],
    'dot-location': ['error', 'property'],
    'dot-notation': ['error'],
    'eol-last': ['error', 'always'],
    'eqeqeq': ['error', 'always', { 'null': 'ignore' }],
    'func-call-spacing': ['error', 'never'],
    'function-paren-newline': ['error', 'multiline'],
    'generator-star-spacing': ['off'], // allow async-await
    'implicit-arrow-linebreak': ['error'],
    'indent': ['error', 2],
    'key-spacing': ['error', { beforeColon: false, afterColon: true, mode: 'strict' }],
    'keyword-spacing': ['error'],
    'linebreak-style': ['error', 'unix'],
    'lines-between-class-members': ['error'],
    'multiline-comment-style': ['warn'],
    'newline-per-chained-call': ['error'],
    'no-console': ['off'],
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off', // allow debugger during development
    'no-else-return': ['error'],
    'no-extra-parens': ['error'],
    'no-implicit-coercion': ['error'],
    'no-lonely-if': ['error'],
    'no-multiple-empty-lines': ['warn', { max: 2, maxEOF: 0, maxBOF: 0 }],
    'no-multi-spaces': ['error'],
    'no-trailing-spaces': ['error'],
    'no-unneeded-ternary': ['error'],
    'no-useless-return': ['error'],
    'no-whitespace-before-property': ['error'],
    'object-curly-newline': ['error', { consistent: true }],
    'object-curly-spacing': ['error', 'always'],
    'object-shorthand': ['error'],
    'padded-blocks': ['error', 'never'],
    'prefer-arrow-callback': ['error'],
    'prefer-const': ['error'],
    'prefer-object-spread': ['error'],
    'prefer-template': ['error'],
    'quote-props': ['error', 'consistent-as-needed', { keywords: true }],
    'quotes': ['error', 'single', { allowTemplateLiterals: true }],
    'semi': ['error', 'never'],
    'space-before-blocks': ['error', 'always'],
    'spaced-comment': ['warn', 'always'],
    'space-infix-ops': ['error'],
    'space-in-parens': ['error', 'never'],
    'space-unary-ops': ['error', { words: true, nonwords: false }],
    'switch-colon-spacing': ['error'],
    'unicode-bom': ['error', 'never'],
    'wrap-iife': ['error'],
    'yoda': ['error'],
  },
}
