
package definitions

func init(){
  add(`Main`, &defMain{})
}

type defMain struct{}

func (*defMain) String() string {
	return `
<interface>
  <object class="GtkWindow" id="mainWindow">
    <property name="window-position">0</property>
    <property name="default-height">600</property>
    <property name="default-width">200</property>
    <property name="title">$title</property>
    <signal name="destroy" handler="on_close_window_signal" />
    <!-- <property name="icon">we dont know how to use it now</property> -->
    <child>
      <object class="GtkBox" id="Vbox">
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <child>
          <object class="GtkMenuBar" id="menubar">
            <child>
              <object class="GtkMenuItem" id="ContactsMenu">
                <property name="label">$contactsMenu</property>
                <child type="submenu">
                  <object class="GtkMenu" id="menu">
                    <child>
                      <object class="GtkMenuItem" id="addMenu">
                        <property name="label">$addMenu</property>
                        <signal name="activate" handler="on_add_contact_window_signal" />
                      </object>
                    </child>
                  </object>
                </child>
              </object>
            </child>
            <child>
              <object class="GtkMenuItem" id="AccountsMenu">
                <property name="label">$accountsMenu</property>
              </object>
            </child>
            <child>
              <object class="GtkMenuItem" id="ViewMenu">
                <property name="label">$viewMenu</property>
                <child type="submenu">
                  <object class="GtkMenu" id="menu2">
                    <child>
                      <object class="GtkCheckMenuItem" id="CheckItemMerge">
                        <property name="label">$checkItemMerge</property>
                        <signal name="toggled" handler="on_toggled_check_Item_Merge_signal" />
                      </object>
                    </child>
                    <child>
                      <object class="GtkCheckMenuItem" id="CheckItemShowOffline">
                        <property name="label" translatable="yes">Show Offline Contacts</property>
                        <signal name="toggled" handler="on_toggled_check_Item_Show_Offline_signal" />
                      </object>
                    </child>
                  </object>
                </child>
              </object>
            </child>
            <child>
              <object class="GtkMenuItem" id="HelpMenu">
                <property name="label">$helpMenu</property>
                <child type="submenu">
                  <object class="GtkMenu" id="menu3">
                    <child>
                      <object class="GtkMenuItem" id="aboutMenu">
                        <property name="label">$aboutMenu</property>
                        <signal name="activate" handler="on_about_dialog_signal" />
                      </object>
                    </child>
                  </object>
                </child>
              </object>
            </child>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">false</property>
            <property name="position">0</property>
          </packing>
        </child>
      </object>
    </child>
  </object>
</interface>

`
}