/* eslint-env node */
import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Components|Gutter Menu',
};

export const GutterMenu = () => {
  return {
    template: hbs`
      <h5 class="title is-5">Gutter Menu</h5>
      <div class="columns">
        <div class="column is-4">
          <div class="gutter">
            <aside class="menu">
              <p class="menu-label">Places</p>
              <ul class="menu-list">
                <li><a href="#" class="is-active">Place One</a></li>
                <li><a href="#">Place Two</a></li>
              </ul>

              <p class="menu-label">Features</p>
              <ul class="menu-list">
                <li><a href="#">Feature One</a></li>
                <li><a href="#">Feature Two</a></li>
              </ul>
            </aside>
          </div>
        </div>
        <div class="column">
          <div class="mock-content">
            <div class="mock-vague"></div>
          </div>
        </div>
      </div>
        `,
  };
};

export const RichComponents = () => {
  return {
    template: hbs`
      <h5 class="title is-5">Gutter Navigation With Rich Components</h5>
      <div class="columns">
        <div class="column is-4">
          <div class="gutter">
            <aside class="menu">
              <p class="menu-label">Places</p>
              <ul class="menu-list">
                <li>
                  <div class="menu-item">
                    <PowerSelect @selected={{or selection "One"}} @options={{array "One" "Two" "Three"}} @onChange={{action (mut selection)}} as |option|>
                      {{option}}
                    </PowerSelect>
                  </div>
                </li>
                <li><a href="#" class="is-active">Place One</a></li>
                <li><a href="#">Place Two</a></li>
              </ul>

              <p class="menu-label">Features</p>
              <ul class="menu-list">
                <li><a href="#">Feature One</a></li>
                <li><a href="#">Feature Two</a></li>
              </ul>
            </aside>
          </div>
        </div>
        <div class="column">
          <div class="mock-content">
            <div class="mock-vague"></div>
          </div>
        </div>
      </div>
      <p class='annotation'>In order to keep the gutter navigation streamlined and easy to navigation, rich components should be avoided when possible. When not possible, they should be kept near the top.</p>
        `,
  };
};

export const ManyItems = () => {
  return {
    template: hbs`
      <h5 class="title is-5">Hypothetical Gutter Navigation With Many Items</h5>
      <div class="columns">
        <div class="column is-4">
          <div class="gutter">
            <aside class="menu">
              <p class="menu-label">Places</p>
              <ul class="menu-list">
                {{#each (array "One Two" "Three" "Four" "Five" "Six" "Seven") as |item|}}
                  <li><a href="#">Place {{item}}</a></li>
                {{/each}}
              </ul>

              <p class="menu-label">Features</p>
              <ul class="menu-list">
                {{#each (array "One Two" "Three" "Four" "Five" "Six" "Seven") as |item|}}
                  <li><a href="#">Feature {{item}}</a></li>
                {{/each}}
              </ul>

              <p class="menu-label">Other</p>
              <ul class="menu-list">
                <li><a href="#" class="is-active">The one that didn't fit in</a></li>
              </ul>

              <p class="menu-label">Things</p>
              <ul class="menu-list">
                {{#each (array "One Two" "Three" "Four" "Five" "Six" "Seven") as |item|}}
                  <li><a href="#">Thing {{item}}</a></li>
                {{/each}}
              </ul>
            </aside>
          </div>
        </div>
        <div class="column">
          <div class="mock-content">
            <div class="mock-vague"></div>
          </div>
        </div>
      </div>
      <p class='annotation'>There will only ever be one gutter menu in the Nomad UI, but it helps to imagine a situation where there are many navigation items in the gutter.</p>
        `,
  };
};

export const IconItems = () => {
  return {
    template: hbs`
      <h5 class="title is-5">Hypothetical Gutter Navigation With Icon Items</h5>
      <div class="columns">
        <div class="column is-4">
          <div class="gutter">
            <aside class="menu">
              <p class="menu-label">Places</p>
              <ul class="menu-list">
                <li><a href="#">{{x-icon "clock"}} Place One</a></li>
                <li><a href="#" class="is-active">{{x-icon "history"}} Place Two</a></li>
              </ul>

              <p class="menu-label">Features</p>
              <ul class="menu-list">
                <li><a href="#">{{x-icon "warning"}} Feature One</a></li>
                <li><a href="#">{{x-icon "media-pause"}} Feature Two</a></li>
              </ul>
            </aside>
          </div>
        </div>
        <div class="column">
          <div class="mock-content">
            <div class="mock-vague"></div>
          </div>
        </div>
      </div>
      <p class='annotation'>In the future, the gutter menu may have icons.</p>
        `,
  };
};

export const Global = () => {
  return {
    template: hbs`
      <h5 class="title is-5">Global Gutter Navigation</h5>
      <div class="columns">
        <div class="column is-4">
          <GutterMenu>
            {{!-- Page content here --}}
          </GutterMenu>
        </div>
      </div>
      <p class='annotation'>Since there will only ever be one gutter menu in the UI, it makes sense to express the menu as a singleton component. This is what that singleton component looks like.</p>
      <p class='annotation'><strong>Note:</strong> Normally the gutter menu is rendered within a page layout and is fixed position. The columns shown in this example are only to imitate the actual width without applying fixed positioning.</p>
        `,
  };
};