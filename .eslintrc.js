module.exports = {
    'env': {
        'commonjs': true,
        'es2021': true,
        'node': true
    },
    'extends': 'eslint:recommended',
    'parserOptions': {
        'ecmaVersion': 13
    },
    'rules': {
        'indent': [
            'warn',
            4
        ],
        'linebreak-style': [
            'warn',
            'windows'
        ],
        'quotes': [
            'warn',
            'single'
        ],
        'semi': [
            'warn',
            'always'
        ]
    }
};
